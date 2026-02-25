// Package authmw provee las primitivas de contexto HTTP para transportar
// la identidad del usuario autenticado a traves de la cadena de middleware.
package authmw

import (
	"net/http"
	"strings"

	"golabs-api/internal/infrastructure/security"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuth es un middleware que valida el JWT en el header Authorization: Bearer.
//
// Flujo:
//  1. Extrae el token del header Authorization.
//  2. Valida la firma y la expiracion usando JWTService.
//  3. Extrae los claims "sub" (userID) y "role" del token.
//  4. Almacena un UserContext en el contexto de la peticion y pasa al siguiente handler.
//
// Rechaza la peticion con HTTP 401 si el header no existe, el token es invalido
// o el claim "sub" esta ausente.
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
				http.Error(w, "token invalido", http.StatusUnauthorized)
				return
			}

			claims := token.Claims.(jwt.MapClaims)
			userID, ok := claims["sub"].(string)
			if !ok || userID == "" {
				http.Error(w, "token invalido", http.StatusUnauthorized)
				return
			}

			role, _ := claims["role"].(string)

			// Almacena el contexto de usuario para los handlers y middleware posteriores.
			ctx := WithUser(r.Context(), UserContext{UserID: userID, Role: role})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
