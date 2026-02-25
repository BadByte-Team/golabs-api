// Package accessmw agrupa los middlewares de control de acceso basado en roles y estado
// del usuario, que se aplican despues de los middlewares de autenticacion (JWTAuth + LoadUser).
package accessmw

import (
	"net/http"

	authmw "golabs-api/internal/interfaces/http/middleware/auth"
)

// RequireRole verifica que el usuario autenticado tenga el rol especificado.
// Retorna HTTP 403 si el rol no coincide.
//
// Uso tipico: proteger rutas de administrador con RequireRole("admin").
func RequireRole(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			user, ok := authmw.GetUser(r.Context())
			if !ok || user.Role != role {
				http.Error(w, "forbidden", http.StatusForbidden)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
