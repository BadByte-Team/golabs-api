package http

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	"golabs-api/internal/health"
	"golabs-api/internal/infrastructure/security"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	userhttp "golabs-api/internal/user/interfaces"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)
	r.Use(authmw.ErrorLogger)

	r.Get("/health", health.NewHandler(db).ServeHTTP)

	jwtSvc, err := security.NewJWTService()
	if err != nil {
		panic("JWT no configurado: " + err.Error())
	}

	// Registrar todas las rutas del dominio user (incluye /auth y /users)
	userhttp.RegisterRoutes(r, db, jwtSvc)

	return r
}
