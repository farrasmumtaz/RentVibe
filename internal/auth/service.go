package auth

import (
	"errors"
	"strings"

	"github.com/farrasmumtaz/RentVibe/internal/models"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var (
	ErrEmailAlreadyUsed  = errors.New("email already used")
	ErrInvalidCredential = errors.New("invalid email or password")
)

type Service interface {
	Register(req RegisterRequest) (*AuthResponse, error)
	Login(req LoginRequest) (*AuthResponse, error)
}

type service struct {
	repository   Repository
	tokenService TokenService
}

func NewService(repository Repository, tokenService TokenService) Service {
	return &service{
		repository:   repository,
		tokenService: tokenService,
	}
}

func (s *service) Register(req RegisterRequest) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	existingUser, err := s.repository.FindByEmail(email)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existingUser != nil {
		return nil, ErrEmailAlreadyUsed
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := models.User{
		Name:         strings.TrimSpace(req.Name),
		Email:        email,
		PasswordHash: string(passwordHash),
	}

	if err := s.repository.Create(&user); err != nil {
		return nil, err
	}

	return s.authResponse(&user)
}

func (s *service) Login(req LoginRequest) (*AuthResponse, error) {
	email := strings.ToLower(strings.TrimSpace(req.Email))

	user, err := s.repository.FindByEmail(email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrInvalidCredential
		}

		return nil, err
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredential
	}

	return s.authResponse(user)
}

func (s *service) authResponse(user *models.User) (*AuthResponse, error) {
	token, err := s.tokenService.Generate(user.ID, user.Email)
	if err != nil {
		return nil, err
	}

	return &AuthResponse{
		Token: token,
		User: UserResponse{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
		},
	}, nil
}

func UserIDFromClaims(claims *Claims) (uint, error) {
	return userIDFromClaims(claims)
}
