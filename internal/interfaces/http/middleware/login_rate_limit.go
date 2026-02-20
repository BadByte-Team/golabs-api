package middleware

import (
	"net"
	"net/http"
	"strings"
	"time"
)

var loginLimiter = newRateLimiter(5, time.Minute)

func LoginRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ip := clientIP(r)
		email := r.URL.Query().Get("email")

		key := ip + ":" + strings.ToLower(email)

		if !loginLimiter.allow(key) {
			http.Error(w, "demasiados intentos, intenta más tarde", http.StatusTooManyRequests)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func clientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
