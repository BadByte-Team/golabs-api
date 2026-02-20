package middleware

import "net/http"

func RequireNotBanned(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		user, ok := GetUser(r.Context())
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
