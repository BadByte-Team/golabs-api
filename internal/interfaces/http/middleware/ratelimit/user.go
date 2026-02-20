package ratelimit

import (
	"net/http"
	"time"

	authmw "golabs-api/internal/interfaces/http/middleware/auth"
)

var userLimiter = newRateLimiter(60, time.Minute)

// UserRateLimit limits requests per authenticated user (60 per minute).
func UserRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := authmw.GetUser(r.Context())
		if !ok || user.UserID == "" {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		if !userLimiter.allow(user.UserID) {
			http.Error(w, "rate limit excedido", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}
