// Package accessmw agrupa los middlewares de control de acceso basado en roles y estado
// del usuario, que se aplican despues de los middlewares de autenticacion (JWTAuth + LoadUser).
package accessmw

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	userdomain "golabs-api/internal/user/domain"
)

// RequireSelfOrAdmin permite el acceso solo al propietario del recurso o a un administrador.
//
// Extrae el param "id" de la URL y lo compara con el UserID del token.
// Un admin puede operar sobre cualquier recurso; un usuario regular solo sobre el suyo.
//
// Retorna HTTP 400 si no hay param "id" en la URL.
// Retorna HTTP 401 si no hay usuario autenticado.
// Retorna HTTP 403 si el usuario intenta acceder a un recurso ajeno sin ser admin.
func RequireSelfOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user, ok := authmw.GetUser(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		targetID := chi.URLParam(r, "id")
		if targetID == "" {
			http.Error(w, "id requerido", http.StatusBadRequest)
			return
		}

		// El admin puede operar sobre cualquier usuario.
		if user.Role == userdomain.RoleAdmin {
			next.ServeHTTP(w, r)
			return
		}

		// El usuario regular solo puede operar sobre su propio recurso.
		if user.UserID != targetID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
