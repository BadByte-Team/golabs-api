package accessmw

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	userdomain "golabs-api/internal/user/domain"
)

// RequireSelfOrAdmin allows access only to the resource owner or an admin.
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

		if user.Role == userdomain.RoleAdmin {
			next.ServeHTTP(w, r)
			return
		}

		if user.UserID != targetID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
