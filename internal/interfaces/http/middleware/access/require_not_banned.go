// Package accessmw agrupa los middlewares de control de acceso basado en roles y estado
// del usuario, que se aplican despues de los middlewares de autenticacion (JWTAuth + LoadUser).
package accessmw

import (
	"net/http"

	authmw "golabs-api/internal/interfaces/http/middleware/auth"
)

// RequireNotBanned rechaza peticiones de usuarios con la cuenta suspendida.
// Retorna HTTP 403 con el mensaje "cuenta suspendida" si el usuario esta baneado.
//
// Debe aplicarse despues de LoadUser, ya que es ese middleware el que popula el campo Banned.
func RequireNotBanned(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := authmw.GetUser(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		if user.Banned {
			http.Error(w, "cuenta suspendida", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}
