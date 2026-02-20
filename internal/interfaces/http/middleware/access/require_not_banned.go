package accessmw

import (
	"net/http"

	authmw "golabs-api/internal/interfaces/http/middleware/auth"
)

// RequireNotBanned rejects requests from banned users.
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
