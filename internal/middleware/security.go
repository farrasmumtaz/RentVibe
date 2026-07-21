package middleware

import (
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

func SecurityHeaders() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
		ctx.Header("X-Content-Type-Options", "nosniff")
		ctx.Header("X-Frame-Options", "DENY")
		ctx.Header("Referrer-Policy", "no-referrer")
		ctx.Header("Content-Security-Policy", "default-src 'none'; frame-ancestors 'none'")

		ctx.Next()
	}
}

func CORS() gin.HandlerFunc {
	allowedOrigins := parseAllowedOrigins(os.Getenv("ALLOWED_ORIGINS"))
	allowedMethods := "GET,POST,PUT,PATCH,DELETE,OPTIONS"
	allowedHeaders := "Authorization,Content-Type"

	return func(ctx *gin.Context) {
		origin := ctx.GetHeader("Origin")
		if origin != "" && isAllowedOrigin(origin, allowedOrigins) {
			ctx.Header("Access-Control-Allow-Origin", origin)
			ctx.Header("Vary", "Origin")
			ctx.Header("Access-Control-Allow-Methods", allowedMethods)
			ctx.Header("Access-Control-Allow-Headers", allowedHeaders)
			ctx.Header("Access-Control-Max-Age", "86400")
		}

		if ctx.Request.Method == http.MethodOptions {
			if origin == "" || !isAllowedOrigin(origin, allowedOrigins) {
				ctx.AbortWithStatus(http.StatusForbidden)
				return
			}

			ctx.AbortWithStatus(http.StatusNoContent)
			return
		}

		ctx.Next()
	}
}

func parseAllowedOrigins(rawOrigins string) map[string]struct{} {
	origins := make(map[string]struct{})

	if rawOrigins == "" {
		rawOrigins = "http://localhost:3000,http://127.0.0.1:3000"
	}

	for _, origin := range strings.Split(rawOrigins, ",") {
		origin = strings.TrimSpace(origin)
		if origin != "" {
			origins[origin] = struct{}{}
		}
	}

	return origins
}

func isAllowedOrigin(origin string, allowedOrigins map[string]struct{}) bool {
	_, exists := allowedOrigins[origin]
	return exists
}
