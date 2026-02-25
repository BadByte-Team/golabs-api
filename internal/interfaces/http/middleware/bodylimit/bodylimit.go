// Package bodylimit provee un middleware para limitar el tamano del body de las peticiones HTTP.
// Previene ataques de agotamiento de recursos enviando bodies extremadamente grandes.
package bodylimit

import (
	"net/http"
)

// MaxBodySize retorna un middleware que limita el body de la peticion a n bytes.
// Las peticiones con body mayor a n bytes reciben HTTP 413 Request Entity Too Large.
//
// Entrada:  n, limite en bytes (p.ej. 1<<20 para 1 MiB).
// Salida:   middleware configurado.
//
// Internamente usa http.MaxBytesReader, que cuando el body supera el limite
// retorna un error al leer y el http.DefaultServeMux responde con 413.
func MaxBodySize(n int64) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			r.Body = http.MaxBytesReader(w, r.Body, n)
			next.ServeHTTP(w, r)
		})
	}
}
