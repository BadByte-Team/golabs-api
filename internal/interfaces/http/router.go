package http

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"golabs-api/internal/health"
	"golabs-api/internal/infrastructure/security"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/middleware/bodylimit"

	challengehttp "golabs-api/internal/challenges/interfaces"
	eventhttp "golabs-api/internal/event/interfaces"
	eventteamhttp "golabs-api/internal/eventteam/interfaces"
	userhttp "golabs-api/internal/user/interfaces"
)

const maxBodyBytes = 1 << 20 // 1 MiB

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// ── Global middleware ──────────────────────────────────────────────────
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(authmw.ErrorLogger)
	r.Use(bodylimit.MaxBodySize(maxBodyBytes))
	r.Use(corsMiddleware())

	// ── Health (no versioning — infra tooling expects fixed paths) ─────────
	h := health.NewHandler(db)
	r.Get("/healthz/live", h.Live)
	r.Get("/healthz/ready", h.Ready)
	// Legacy path kept for existing tooling
	r.Get("/health", h.ServeHTTP)

	jwtSvc, err := security.NewJWTService()
	if err != nil {
		panic("JWT no configurado: " + err.Error())
	}

	// ── Versioned API ──────────────────────────────────────────────────────
	r.Route("/api/v1", func(r chi.Router) {
		userhttp.RegisterRoutes(r, db, jwtSvc)
		eventhttp.RegisterRoutes(r, db, jwtSvc)
		eventteamhttp.RegisterRoutes(r, db, jwtSvc)
		challengehttp.RegisterRoutes(r, db, jwtSvc)
	})

	return r
}

// corsMiddleware builds a CORS handler from the ALLOWED_ORIGINS env var.
// Set ALLOWED_ORIGINS="http://localhost:3000,https://app.example.com"
// Defaults to "*" for development if the variable is not set.
func corsMiddleware() func(http.Handler) http.Handler {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if originsEnv == "" {
		origins = []string{"*"}
	} else {
		for _, o := range strings.Split(originsEnv, ",") {
			if s := strings.TrimSpace(o); s != "" {
				origins = append(origins, s)
			}
		}
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: originsEnv != "", // only when explicit origins are set
		MaxAge:           300,
	})
}
