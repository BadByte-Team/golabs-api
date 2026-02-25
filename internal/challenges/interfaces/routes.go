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

// RegisterRoutes wires all challenge routes under /events/{event_id}/challenges.
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {
	challengeRepo := challengeinfra.NewChallengeRepository(db)
	eventRepo := eventinfra.NewEventRepository(db)
	teamRepo := teaminfra.NewEventTeamRepository(db)
	userRepo := userinfra.NewUserRepository(db)

	// Use cases
	createUC := challengeapp.NewCreateChallengeUseCase(challengeRepo, eventRepo)
	updateUC := challengeapp.NewUpdateChallengeUseCase(challengeRepo)
	publishUC := challengeapp.NewPublishChallengeUseCase(challengeRepo)
	listUC := challengeapp.NewListChallengesUseCase(challengeRepo)
	getUC := challengeapp.NewGetChallengeUseCase(challengeRepo)
	setFlagUC := challengeapp.NewSetFlagUseCase(challengeRepo)
	submitUC := challengeapp.NewSubmitFlagUseCase(challengeRepo, eventRepo, teamRepo)

	h := NewChallengeHandler(createUC, updateUC, publishUC, listUC, getUC, setFlagUC, submitUC)

	r.Route("/events/{event_id}/challenges", func(r chi.Router) {

		// ── Public (authenticated, role already in JWT context) ──────────────
		r.Group(func(r chi.Router) {
			r.Use(authmw.JWTAuth(jwtSvc))

			r.Get("/", h.List)
			r.Get("/{challenge_id}", h.Get)

			// Submit: authenticated + not banned + rate limited
			r.Group(func(r chi.Router) {
				r.Use(authmw.LoadUser(userRepo))
				r.Use(accessmw.RequireNotBanned)
				r.Use(ratelimit.UserRateLimit)
				r.Post("/{challenge_id}/submit", h.Submit)
			})
		})

		// ── Admin only (role from JWT, no extra DB round-trip) ───────────────
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
