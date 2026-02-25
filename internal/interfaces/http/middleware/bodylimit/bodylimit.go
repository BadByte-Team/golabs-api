package bodylimit

import (
	"net/http"
)

// MaxBodySize returns a middleware that limits the request body to n bytes.
// Requests that exceed the limit receive 413 Request Entity Too Large.
func MaxBodySize(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}
