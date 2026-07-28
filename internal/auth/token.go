package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

const defaultTokenTTL = 24 * time.Hour

var (
	ErrInvalidToken     = errors.New("invalid token")
	ErrExpiredToken     = errors.New("expired token")
	ErrMissingJWTSecret = errors.New("JWT_SECRET must be configured")
)

type Claims struct {
	Subject  string `json:"sub"`
	Email    string `json:"email"`
	IssuedAt int64  `json:"iat"`
	Expires  int64  `json:"exp"`
}

type tokenHeader struct {
	Algorithm string `json:"alg"`
	Type      string `json:"typ"`
}

type TokenService interface {
	Generate(userID uint, email string) (string, error)
	Validate(token string) (*Claims, error)
}

type hmacTokenService struct {
	secret []byte
	ttl    time.Duration
}

func NewTokenService() (TokenService, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, ErrMissingJWTSecret
	}

	return &hmacTokenService{
		secret: []byte(secret),
		ttl:    tokenTTLFromEnv(),
	}, nil
}

func (s *hmacTokenService) Generate(userID uint, email string) (string, error) {
	now := time.Now().UTC()
	claims := Claims{
		Subject:  strconv.FormatUint(uint64(userID), 10),
		Email:    email,
		IssuedAt: now.Unix(),
		Expires:  now.Add(s.ttl).Unix(),
	}

	headerPart, err := encodeJSON(tokenHeader{Algorithm: "HS256", Type: "JWT"})
	if err != nil {
		return "", err
	}

	payloadPart, err := encodeJSON(claims)
	if err != nil {
		return "", err
	}

	unsignedToken := headerPart + "." + payloadPart
	signature := s.sign(unsignedToken)

	return unsignedToken + "." + signature, nil
}

func (s *hmacTokenService) Validate(token string) (*Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 {
		return nil, ErrInvalidToken
	}

	header, err := decodeHeader(parts[0])
	if err != nil || header.Algorithm != "HS256" || header.Type != "JWT" {
		return nil, ErrInvalidToken
	}

	unsignedToken := parts[0] + "." + parts[1]
	expectedSignature := s.sign(unsignedToken)
	if !hmac.Equal([]byte(expectedSignature), []byte(parts[2])) {
		return nil, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return nil, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return nil, ErrInvalidToken
	}

	if claims.Expires <= time.Now().UTC().Unix() {
		return nil, ErrExpiredToken
	}

	return &claims, nil
}

func (s *hmacTokenService) sign(unsignedToken string) string {
	hash := hmac.New(sha256.New, s.secret)
	hash.Write([]byte(unsignedToken))
	return base64.RawURLEncoding.EncodeToString(hash.Sum(nil))
}

func encodeJSON(value interface{}) (string, error) {
	payload, err := json.Marshal(value)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(payload), nil
}

func decodeHeader(encodedHeader string) (*tokenHeader, error) {
	payload, err := base64.RawURLEncoding.DecodeString(encodedHeader)
	if err != nil {
		return nil, err
	}

	var header tokenHeader
	if err := json.Unmarshal(payload, &header); err != nil {
		return nil, err
	}

	return &header, nil
}

func tokenTTLFromEnv() time.Duration {
	rawTTL := os.Getenv("JWT_TTL_MINUTES")
	if rawTTL == "" {
		return defaultTokenTTL
	}

	minutes, err := strconv.Atoi(rawTTL)
	if err != nil || minutes <= 0 {
		return defaultTokenTTL
	}

	return time.Duration(minutes) * time.Minute
}

func userIDFromClaims(claims *Claims) (uint, error) {
	id, err := strconv.ParseUint(claims.Subject, 10, 64)
	if err != nil {
		return 0, fmt.Errorf("invalid subject claim: %w", err)
	}

	return uint(id), nil
}
