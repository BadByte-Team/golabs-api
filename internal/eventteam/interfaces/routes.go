package interfaces

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	eventinfra "golabs-api/internal/event/infraestructure"
	"golabs-api/internal/infrastructure/security"
	userinfra "golabs-api/internal/user/infrastructure"

	teamapp "golabs-api/internal/eventteam/application"
	teaminfra "golabs-api/internal/eventteam/infraestructure"

	accessmw "golabs-api/internal/interfaces/http/middleware/access"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
)

// RegisterRoutes wires all event team routes.
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {

	// Repos
	eventRepo := eventinfra.NewEventRepository(db)
	teamRepo := teaminfra.NewEventTeamRepository(db)
	userRepo := userinfra.NewUserRepository(db)

	// Use cases
	createUC := teamapp.NewCreateTeamUseCase(eventRepo, teamRepo)
	joinUC := teamapp.NewJoinTeamUseCase(eventRepo, teamRepo)
	leaveUC := teamapp.NewLeaveTeamUseCase(teamRepo)
	rotateUC := teamapp.NewRotateJoinSecretUseCase(teamRepo)
	listTeamsUC := teamapp.NewListTeamsByEventUseCase(teamRepo)
	leaderboardUC := teamapp.NewGetLeaderboardUseCase(teamRepo)

	// Handler
	handler := NewEventTeamHandler(
		createUC,
		joinUC,
		leaveUC,
		rotateUC,
		listTeamsUC,
		leaderboardUC,
	)

	// Protected routes (user autenticado)
	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuth(jwtSvc))
		r.Use(authmw.LoadUser(userRepo))
		r.Use(accessmw.RequireNotBanned)

		r.Route("/events/{event_id}/teams", func(r chi.Router) {
			// List all teams in an event
			r.Get("/", handler.ListTeams)

			// Create a team
			r.Post("/", handler.Create)

			// Join a team with secret
			r.Post("/join", handler.Join)

			// Actions on a specific team
			r.Route("/{team_id}", func(r chi.Router) {
				r.Get("/members", handler.ListMembers)
				r.Post("/leave", handler.Leave)
				r.Post("/rotate-secret", handler.RotateSecret)
			})
		})

		// Leaderboard is public to any authenticated user
		r.Get("/events/{event_id}/leaderboard", handler.Leaderboard)
	})
}
