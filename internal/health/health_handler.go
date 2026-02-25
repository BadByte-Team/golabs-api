package health

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"time"
)

type Handler struct {
	db *sql.DB
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{db: db}
}

// Live is always 200 — proves the process is running.
// Used by container orchestrators for liveness probes.
func (h *Handler) Live(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

// Ready checks downstream dependencies (DB).
// Returns 200 when ready, 503 when the DB is unavailable.
// Used by orchestrators to gate traffic: pod receives requests only when ready.
func (h *Handler) Ready(w http.ResponseWriter, r *http.Request) {
	dbStatus := "ok"
	code := http.StatusOK

	ctx := r.Context()
	if err := h.db.PingContext(ctx); err != nil {
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

// ServeHTTP kept for backward compatibility with any existing /health usage.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.Ready(w, r)
}
