// Package authmw provee las primitivas de contexto HTTP para transportar
// la identidad del usuario autenticado a traves de la cadena de middleware.
package authmw

import (
	"log/slog"
	"net/http"
)

// responseWriter es un wrapper de http.ResponseWriter que captura el codigo HTTP
// escrito al response, permitiendo que ErrorLogger lo inspeccione despues del handler.
type responseWriter struct {
	http.ResponseWriter
	status int // codigo HTTP real enviado al cliente
}

// WriteHeader captura el codigo de estado antes de delegar al ResponseWriter original.
func (rw *responseWriter) WriteHeader(code int) {
	rw.status = code
	rw.ResponseWriter.WriteHeader(code)
}

// ErrorLogger registra con slog las respuestas HTTP con codigo 5xx (errores de servidor).
// Los errores 4xx (errores de cliente) no se loguean para evitar spam en logs de produccion.
//
// Se posiciona como el primer middleware de la cadena para capturar errores de todos
// los handlers, incluyendo panics recuperados por chi.Recoverer.
func ErrorLogger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rw := &responseWriter{ResponseWriter: w, status: http.StatusOK}

		next.ServeHTTP(rw, r)

		// Loguear solo errores de servidor (5xx) en structured logging.
		if rw.status >= 500 {
			slog.Error(
				"http server error",
				"method", r.Method,
				"path", r.URL.Path,
				"status", rw.status,
				"remote", r.RemoteAddr,
			)
		}
	})
}
