package middleware

import (
	"net/http"
	"time"
)

var userLimiter = newRateLimiter(60, time.Minute)

func UserRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, ok := GetUser(r.Context())
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
