package middleware

import (
	"net/http"
	"strings"

	"golabs-api/internal/infrastructure/security"

	"github.com/golang-jwt/jwt/v5"
)

func JWTAuth(jwtSvc *security.JWTService) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				http.Error(w, "token requerido", http.StatusUnauthorized)
				return
			}

			tokenStr := strings.TrimPrefix(auth, "Bearer ")

			token, err := jwtSvc.Parse(tokenStr)
			if err != nil || !token.Valid {
				http.Error(w, "token inválido", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)

			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				http.Error(w, "token inválido", http.StatusUnauthorized)
				return
			}

			// Identidad
			ctx := WithUser(r.Context(), UserContext{
				UserID: userID,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
