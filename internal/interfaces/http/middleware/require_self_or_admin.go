package middleware

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func RequireSelfOrAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, ok := GetUser(r.Context())
		if !ok {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}

		targetID := chi.URLParam(r, "id")
		if targetID == "" {
			http.Error(w, "id requerido", http.StatusBadRequest)
			return
		}

		// Admin puede todo
		if user.Role == "admin" {
			next.ServeHTTP(w, r)
			return
		}

		// Usuario solo sobre sí mismo
		if user.UserID != targetID {
			http.Error(w, "forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}
