// Package health implementa los endpoints de health check del servidor.
// Sigue el patron de separar liveness y readiness para orquestadores de contenedores
// (Kubernetes, Docker Swarm, etc.) que usan ambas probes con comportamientos distintos.
package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

// Handler agrupa los endpoints de health check y mantiene referencia al pool de BD
// para verificar la disponibilidad de la base de datos en el probe de readiness.
type Handler struct {
	db *sql.DB
}

// NewHandler crea un Handler de health check con acceso al pool de conexiones.
//
// Entrada:  db, pool de conexiones SQL activo.
// Salida:   puntero al Handler configurado.
func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// Live responde siempre con HTTP 200 mientras el proceso este en ejecucion.
// No verifica dependencias externas; su unico proposito es confirmar que el
// proceso esta vivo y puede recibir trafico para ser reiniciado si no responde.
//
// GET /healthz/live
func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Ready verifica que las dependencias criticas (base de datos) esten disponibles.
// Retorna HTTP 200 cuando el servidor puede procesar peticiones de negocio,
// o HTTP 503 cuando la base de datos no responde.
//
// Los orquestadores usan esta probe para decidir si enviar trafico al pod.
//
// GET /healthz/ready
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	code := http.StatusOK

	// Ping con contexto de la peticion para respetar timeouts ya configurados.
	if err := h.db.PingContext(r.Context()); err != nil {
		dbStatus = "unavailable"
		code = http.StatusServiceUnavailable
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":   dbStatus,
		"database": dbStatus,
		"time":     time.Now().UTC(),
	})
}

// ServeHTTP delega a Ready para mantener compatibilidad con el endpoint /health legacy.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Ready(w, r)
}
