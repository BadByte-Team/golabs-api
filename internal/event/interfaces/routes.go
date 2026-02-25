package interfaces

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	eventapp "golabs-api/internal/event/application"
	eventinfra "golabs-api/internal/event/infraestructure"
	"golabs-api/internal/infrastructure/security"
	userdomain "golabs-api/internal/user/domain"

	accessmw "golabs-api/internal/interfaces/http/middleware/access"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
)

// RegisterRoutes wires all event routes.
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {
	repo := eventinfra.NewEventRepository(db)

	// Use cases
	createUC := eventapp.NewCreateEventUseCase(repo)
	getUC := eventapp.NewGetEventByIDUseCase(repo)
	listUC := eventapp.NewListEventsUseCase(repo)
	openUC := eventapp.NewOpenEventUseCase(repo)
	startUC := eventapp.NewStartEventUseCase(repo)
	finishUC := eventapp.NewFinishEventUseCase(repo)

	// Handler
	handler := NewEventHandler(
		createUC,
		getUC,
		listUC,
		openUC,
		startUC,
		finishUC,
	)

	r.Route("/events", func(r chi.Router) {

		r.Get("/", handler.List)
		r.Get("/{event_id}", handler.GetByID)

		// Admin-only: JWTAuth already puts the role from the token into context,
		// so we don't need LoadUser here — saving a DB round-trip per request.
		r.Group(func(r chi.Router) {
			r.Use(authmw.JWTAuth(jwtSvc))
			r.Use(accessmw.RequireRole(userdomain.RoleAdmin))

			r.Post("/", handler.Create)
			r.Post("/{event_id}/open", handler.Open)
			r.Post("/{event_id}/start", handler.Start)
			r.Post("/{event_id}/finish", handler.Finish)
		})
	})
}
