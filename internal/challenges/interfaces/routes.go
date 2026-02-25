// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de challenges.
package interfaces

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	challengeapp "golabs-api/internal/challenges/application"
	challengeinfra "golabs-api/internal/challenges/infraestructure"
	eventinfra "golabs-api/internal/event/infraestructure"
	teaminfra "golabs-api/internal/eventteam/infraestructure"
	"golabs-api/internal/infrastructure/security"
	accessmw "golabs-api/internal/interfaces/http/middleware/access"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/middleware/ratelimit"
	userdomain "golabs-api/internal/user/domain"
	userinfra "golabs-api/internal/user/infrastructure"
)

// RegisterRoutes registra todas las rutas del modulo de challenges bajo /events/{event_id}/challenges.
//
// Rutas publicas (solo JWT, sin LoadUser para evitar consulta extra):
//   - GET  /events/{event_id}/challenges                      lista de challenges del evento
//   - GET  /events/{event_id}/challenges/{challenge_id}       detalle de challenge
//
// Ruta de participantes (JWT + LoadUser + no baneado + rate limit):
//   - POST /events/{event_id}/challenges/{challenge_id}/submit  enviar flag
//
// Rutas de admin (JWT + rol admin, sin LoadUser):
//   - POST /events/{event_id}/challenges                        crear challenge
//   - PUT  /events/{event_id}/challenges/{challenge_id}         actualizar challenge
//   - POST /events/{event_id}/challenges/{challenge_id}/publish  publicar
//   - POST /events/{event_id}/challenges/{challenge_id}/unpublish ocultar
//   - POST /events/{event_id}/challenges/{challenge_id}/flag    establecer flag
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {
	challengeRepo := challengeinfra.NewChallengeRepository(db)
	eventRepo := eventinfra.NewEventRepository(db)
	teamRepo := teaminfra.NewEventTeamRepository(db)
	userRepo := userinfra.NewUserRepository(db)

	// Use cases.
	createUC := challengeapp.NewCreateChallengeUseCase(challengeRepo, eventRepo)
	updateUC := challengeapp.NewUpdateChallengeUseCase(challengeRepo)
	publishUC := challengeapp.NewPublishChallengeUseCase(challengeRepo)
	listUC := challengeapp.NewListChallengesUseCase(challengeRepo)
	getUC := challengeapp.NewGetChallengeUseCase(challengeRepo)
	setFlagUC := challengeapp.NewSetFlagUseCase(challengeRepo)
	submitUC := challengeapp.NewSubmitFlagUseCase(challengeRepo, eventRepo, teamRepo)

	h := NewChallengeHandler(createUC, updateUC, publishUC, listUC, getUC, setFlagUC, submitUC)

	r.Route("/events/{event_id}/challenges", func(r chi.Router) {

		// Rutas publicas para usuarios autenticados (el rol viene del JWT, sin LoadUser).
		r.Group(func(r chi.Router) {
			r.Use(authmw.JWTAuth(jwtSvc))

			r.Get("/", h.List)
			r.Get("/{challenge_id}", h.Get)

			// Submit: requiere a su vez que el usuario este activo (no baneado).
			r.Group(func(r chi.Router) {
				r.Use(authmw.LoadUser(userRepo))
				r.Use(accessmw.RequireNotBanned)
				r.Use(ratelimit.UserRateLimit)
				r.Post("/{challenge_id}/submit", h.Submit)
			})
		})

		// Rutas de administracion (solo admin, sin LoadUser: el rol viene del JWT).
		r.Group(func(r chi.Router) {
			r.Use(authmw.JWTAuth(jwtSvc))
			r.Use(accessmw.RequireRole(userdomain.RoleAdmin))

			r.Post("/", h.Create)
			r.Put("/{challenge_id}", h.Update)
			r.Post("/{challenge_id}/publish", h.Publish)
			r.Post("/{challenge_id}/unpublish", h.Unpublish)
			r.Post("/{challenge_id}/flag", h.SetFlag)
		})
	})
}
