// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de eventos.
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

// RegisterRoutes registra todas las rutas del modulo de eventos en el router dado.
//
// Rutas publicas (sin autenticacion):
//   - GET /events/        lista de eventos (paginada)
//   - GET /events/{id}    detalle de un evento
//
// Rutas de admin (requieren JWT + rol "admin"):
//   - POST   /events/              crear evento (estado inicial: draft)
//   - PUT    /events/{id}          actualizar evento (solo draft)
//   - DELETE /events/{id}          eliminar evento (solo draft)
//   - POST   /events/{id}/open     transicion draft -> open
//   - POST   /events/{id}/start    transicion open  -> running
//   - POST   /events/{id}/finish   transicion running -> finished
//
// Nota: las rutas de admin usan JWTAuth sin LoadUser ya que el rol viene del token,
// evitando una consulta adicional a BD por peticion.
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {
	repo := eventinfra.NewEventRepository(db)

	// Instanciar use cases.
	createUC := eventapp.NewCreateEventUseCase(repo)
	updateUC := eventapp.NewUpdateEventUseCase(repo)
	deleteUC := eventapp.NewDeleteEventUseCase(repo)
	getUC := eventapp.NewGetEventByIDUseCase(repo)
	listUC := eventapp.NewListEventsUseCase(repo)
	openUC := eventapp.NewOpenEventUseCase(repo)
	startUC := eventapp.NewStartEventUseCase(repo)
	finishUC := eventapp.NewFinishEventUseCase(repo)

	// Instanciar handler.
	handler := NewEventHandler(
		createUC,
		updateUC,
		deleteUC,
		getUC,
		listUC,
		openUC,
		startUC,
		finishUC,
	)

	r.Route("/events", func(r chi.Router) {

		// Rutas publicas: no requieren autenticacion.
		r.Get("/", handler.List)
		r.Get("/{event_id}", handler.GetByID)

		// Rutas de admin: JWTAuth es suficiente (el rol viene del token, sin LoadUser).
		r.Group(func(r chi.Router) {
			r.Use(authmw.JWTAuth(jwtSvc))
			r.Use(accessmw.RequireRole(userdomain.RoleAdmin))

			r.Post("/", handler.Create)
			r.Put("/{event_id}", handler.Update)
			r.Post("/{event_id}/delete", handler.Delete)
			r.Post("/{event_id}/open", handler.Open)
			r.Post("/{event_id}/start", handler.Start)
			r.Post("/{event_id}/finish", handler.Finish)
		})
	})
}
