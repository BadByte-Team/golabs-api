package http

import (
	"net/http"

	"golabs-api/internal/health"
)

func NewRouter() http.Handler {
	mux := http.NewServeMux()

	healthHandler := health.NewHandler()

	mux.Handle("/health", healthHandler)

	return mux
}
