package middleware

import (
	"errors"
	"net/http"
	"strings"

	"github.com/farrasmumtaz/RentVibe/internal/auth"
	"github.com/farrasmumtaz/RentVibe/pkg/response"

	"github.com/gin-gonic/gin"
)

const (
	authorizationHeader = "Authorization"
	bearerPrefix        = "Bearer "
	userIDContextKey    = "user_id"
	userEmailContextKey = "user_email"
)

func Auth(tokenService auth.TokenService) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		header := ctx.GetHeader(authorizationHeader)
		if !strings.HasPrefix(header, bearerPrefix) {
			response.Error(ctx, http.StatusUnauthorized, "Token otorisasi wajib menggunakan format Bearer")
			ctx.Abort()
			return
		}

		token := strings.TrimSpace(strings.TrimPrefix(header, bearerPrefix))
		if token == "" {
			response.Error(ctx, http.StatusUnauthorized, "Token otorisasi tidak ditemukan")
			ctx.Abort()
			return
		}

		claims, err := tokenService.Validate(token)
		if err != nil {
			message := "Token otorisasi tidak valid"
			if errors.Is(err, auth.ErrExpiredToken) {
				message = "Token otorisasi sudah kadaluarsa"
			}

			response.Error(ctx, http.StatusUnauthorized, message)
			ctx.Abort()
			return
		}

		userID, err := auth.UserIDFromClaims(claims)
		if err != nil {
			response.Error(ctx, http.StatusUnauthorized, "Token otorisasi tidak valid")
			ctx.Abort()
			return
		}

		ctx.Set(userIDContextKey, userID)
		ctx.Set(userEmailContextKey, claims.Email)
		ctx.Next()
	}
}
