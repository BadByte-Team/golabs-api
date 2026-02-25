// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de equipos.
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

// RegisterRoutes registra todas las rutas del modulo de equipos en el router dado.
//
// Todas las rutas requieren autenticacion (JWTAuth + LoadUser + no baneado).
//
// Rutas de equipo:
//   - GET  /events/{event_id}/teams                      lista equipos del evento
//   - POST /events/{event_id}/teams                      crear equipo (el creador es owner)
//   - POST /events/{event_id}/teams/join                 unirse a equipo con join secret
//   - GET  /events/{event_id}/teams/{team_id}/members    listar miembros del equipo
//   - POST /events/{event_id}/teams/{team_id}/leave      abandonar equipo
//   - POST /events/{event_id}/teams/{team_id}/rotate-secret  rotar join secret (solo owner)
//
// Ruta de leaderboard:
//   - GET  /events/{event_id}/leaderboard    ranking de equipos por puntaje
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {

	// Repositorios.
	eventRepo := eventinfra.NewEventRepository(db)
	teamRepo := teaminfra.NewEventTeamRepository(db)
	userRepo := userinfra.NewUserRepository(db)

	// Use cases.
	createUC := teamapp.NewCreateTeamUseCase(eventRepo, teamRepo)
	joinUC := teamapp.NewJoinTeamUseCase(eventRepo, teamRepo)
	leaveUC := teamapp.NewLeaveTeamUseCase(teamRepo)
	rotateUC := teamapp.NewRotateJoinSecretUseCase(teamRepo)
	listTeamsUC := teamapp.NewListTeamsByEventUseCase(teamRepo)
	leaderboardUC := teamapp.NewGetLeaderboardUseCase(teamRepo)

	// Handler.
	handler := NewEventTeamHandler(
		createUC,
		joinUC,
		leaveUC,
		rotateUC,
		listTeamsUC,
		leaderboardUC,
	)

	// Todas las rutas requieren autenticacion y usuario activo (no baneado).
	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuth(jwtSvc))
		r.Use(authmw.LoadUser(userRepo))
		r.Use(accessmw.RequireNotBanned)

		r.Route("/events/{event_id}/teams", func(r chi.Router) {
			r.Get("/", handler.ListTeams)
			r.Post("/", handler.Create)
			r.Post("/join", handler.Join)

			r.Route("/{team_id}", func(r chi.Router) {
				r.Get("/members", handler.ListMembers)
				r.Post("/leave", handler.Leave)
				r.Post("/rotate-secret", handler.RotateSecret)
			})
		})

		// El leaderboard es accesible para cualquier usuario autenticado.
		r.Get("/events/{event_id}/leaderboard", handler.Leaderboard)
	})
}
