package middleware

import (
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/farrasmumtaz/RentVibe/pkg/response"

	"github.com/gin-gonic/gin"
)

type visitor struct {
	count     int
	expiresAt time.Time
}

func RateLimiter(maxRequests int, window time.Duration) gin.HandlerFunc {
	visitors := make(map[string]*visitor)
	var mu sync.Mutex

	go func() {
		ticker := time.NewTicker(window)
		defer ticker.Stop()

		for range ticker.C {
			now := time.Now()
			mu.Lock()
			for ip, currentVisitor := range visitors {
				if now.After(currentVisitor.expiresAt) {
					delete(visitors, ip)
				}
			}
			mu.Unlock()
		}
	}()

	return func(ctx *gin.Context) {
		ip := clientIP(ctx)
		now := time.Now()

		mu.Lock()
		currentVisitor, exists := visitors[ip]
		if !exists || now.After(currentVisitor.expiresAt) {
			currentVisitor = &visitor{
				count:     0,
				expiresAt: now.Add(window),
			}
			visitors[ip] = currentVisitor
		}

		currentVisitor.count++
		retryAfter := time.Until(currentVisitor.expiresAt)
		allowed := currentVisitor.count <= maxRequests
		mu.Unlock()

		if !allowed {
			ctx.Header("Retry-After", strconv.Itoa(retryAfterSeconds(retryAfter)))
			response.Error(ctx, http.StatusTooManyRequests, "Terlalu banyak request, coba lagi nanti")
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}

func retryAfterSeconds(duration time.Duration) int {
	if duration <= 0 {
		return 1
	}

	return int((duration + time.Second - 1) / time.Second)
}

func clientIP(ctx *gin.Context) string {
	ip := ctx.ClientIP()
	if parsedIP := net.ParseIP(ip); parsedIP != nil {
		return parsedIP.String()
	}

	return ip
}
