package middleware

import "net/http"

// ⚠️ TEMPORAL — luego JWT
func FakeAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// Simulación:
		user := UserContext{
			UserID: r.Header.Get("X-User-ID"),
			Role:   r.Header.Get("X-User-Role"), // "user" | "admin"
		}

		ctx := WithUser(r.Context(), user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
