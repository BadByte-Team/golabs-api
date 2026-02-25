// Package ratelimit implementa un limitador de tasa en memoria basado en una ventana deslizante.
package ratelimit

import (
	"net/http"
	"time"

	authmw "golabs-api/internal/interfaces/http/middleware/auth"
)

// userLimiter permite como maximo 60 requests por usuario autenticado por minuto.
// Protege los endpoints internos contra uso abusivo por parte de usuarios validos.
var userLimiter = newRateLimiter(60, time.Minute)

// UserRateLimit limita los requests por ID de usuario autenticado (60 por minuto).
//
// Requiere que JWTAuth haya corrido antes para que el UserContext este en el contexto.
// Retorna HTTP 401 si no hay usuario autenticado.
// Retorna HTTP 429 Too Many Requests si el usuario supera el limite.
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
