// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de equipos.
package interfaces

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"golabs-api/internal/apperrors"
	eventteamapp "golabs-api/internal/eventteam/application"
	teamdomain "golabs-api/internal/eventteam/domain"
	authctx "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/validate"
)

// EventTeamHandler agrupa los handlers HTTP del modulo de equipos por evento.
type EventTeamHandler struct {
	createUC      *eventteamapp.CreateTeamUseCase
	joinUC        *eventteamapp.JoinTeamUseCase
	leaveUC       *eventteamapp.LeaveTeamUseCase
	rotateUC      *eventteamapp.RotateJoinSecretUseCase
	listTeamsUC   *eventteamapp.ListTeamsByEventUseCase
	leaderboardUC *eventteamapp.GetLeaderboardUseCase
}

// NewEventTeamHandler inyecta las dependencias del EventTeamHandler.
func NewEventTeamHandler(
	create *eventteamapp.CreateTeamUseCase,
	join *eventteamapp.JoinTeamUseCase,
	leave *eventteamapp.LeaveTeamUseCase,
	rotate *eventteamapp.RotateJoinSecretUseCase,
	listTeams *eventteamapp.ListTeamsByEventUseCase,
	leaderboard *eventteamapp.GetLeaderboardUseCase,
) *EventTeamHandler {
	return &EventTeamHandler{
		createUC:      create,
		joinUC:        join,
		leaveUC:       leave,
		rotateUC:      rotate,
		listTeamsUC:   listTeams,
		leaderboardUC: leaderboard,
	}
}

// Create godoc — POST /api/v1/events/{event_id}/teams
// Crea un nuevo equipo en el evento. El usuario autenticado queda como owner.
// Body:  CreateTeamRequest | Exito: 201 CreateTeamResponse (incluye join_secret UNA SOLA VEZ)
func (h *EventTeamHandler) Create(w http.ResponseWriter, r *http.Request) {
	user, _ := authctx.GetUser(r.Context())

	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	userID, err := uuid.Parse(user.UserID)
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrUnauthorized)
		return
	}

	var req CreateTeamRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.createUC.Execute(eventID, userID, req.Name)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusCreated, CreateTeamResponse{
		EventTeamResponse: EventTeamResponse{
			ID:      result.Team.ID.String(),
			EventID: result.Team.EventID.String(),
			Name:    result.Team.Name,
			Score:   result.Team.Score,
		},
		JoinSecret: result.JoinSecret,
	})
}

// Join godoc — POST /api/v1/events/{event_id}/teams/join
// Une al usuario autenticado a un equipo existente usando el join secret.
// Body:  JoinTeamRequest | Exito: 204 No Content
func (h *EventTeamHandler) Join(w http.ResponseWriter, r *http.Request) {
	user, _ := authctx.GetUser(r.Context())

	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	userID, err := uuid.Parse(user.UserID)
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrUnauthorized)
		return
	}

	var req JoinTeamRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	if err = h.joinUC.Execute(eventID, userID, req.TeamName, req.JoinSecret); err != nil {
		apperrors.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// Leave godoc — DELETE /api/v1/events/{event_id}/teams/{team_id}/leave
// El usuario autenticado abandona el equipo. El owner solo puede salir si es el unico miembro.
// Exito: 204 No Content
func (h *EventTeamHandler) Leave(w http.ResponseWriter, r *http.Request) {
	user, _ := authctx.GetUser(r.Context())

	teamID, err := uuid.Parse(chi.URLParam(r, "team_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	userID, err := uuid.Parse(user.UserID)
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrUnauthorized)
		return
	}

	if err := h.leaveUC.Execute(teamID, userID); err != nil {
		apperrors.RespondError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// RotateSecret godoc — POST /api/v1/events/{event_id}/teams/{team_id}/rotate-secret
// El owner rota el join secret del equipo, invalidando el anterior.
// Exito: 200 { "join_secret": "<nuevo_secret>" }
func (h *EventTeamHandler) RotateSecret(w http.ResponseWriter, r *http.Request) {
	user, _ := authctx.GetUser(r.Context())

	teamID, err := uuid.Parse(chi.URLParam(r, "team_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	userID, err := uuid.Parse(user.UserID)
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrUnauthorized)
		return
	}

	secret, err := h.rotateUC.Execute(teamID, userID)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, map[string]string{"join_secret": secret})
}

// ListTeams godoc — GET /api/v1/events/{event_id}/teams
// Retorna todos los equipos inscritos en el evento.
// Exito: 200 []EventTeamResponse
func (h *EventTeamHandler) ListTeams(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	teams, err := h.listTeamsUC.Execute(eventID)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	resp := make([]EventTeamResponse, 0, len(teams))
	for _, t := range teams {
		resp = append(resp, EventTeamResponse{
			ID:      t.ID.String(),
			EventID: t.EventID.String(),
			Name:    t.Name,
			Score:   t.Score,
		})
	}
	apperrors.RespondJSON(w, http.StatusOK, resp)
}

// ListMembers godoc — GET /api/v1/events/{event_id}/teams/{team_id}/members
// Retorna todos los miembros del equipo con sus usernames y roles.
// Exito: 200 []EventTeamMemberResponse
func (h *EventTeamHandler) ListMembers(w http.ResponseWriter, r *http.Request) {
	teamID, err := uuid.Parse(chi.URLParam(r, "team_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	members, err := h.listTeamsUC.ExecuteMembers(teamID)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	resp := make([]EventTeamMemberResponse, 0, len(members))
	for _, m := range members {
		resp = append(resp, EventTeamMemberResponse{
			UserID:   m.UserID.String(),
			Username: m.Username,
			Role:     string(m.Role),
			JoinedAt: m.JoinedAt,
		})
	}
	apperrors.RespondJSON(w, http.StatusOK, resp)
}

// Leaderboard godoc — GET /api/v1/events/{event_id}/leaderboard
// Retorna el ranking de equipos del evento ordenado por puntaje descendente.
// Exito: 200 []LeaderboardEntry
func (h *EventTeamHandler) Leaderboard(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	entries, err := h.leaderboardUC.Execute(eventID)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, entries)
}

// ── helpers ───────────────────────────────────────────────────────────────────

// MemberWithUsername es un alias del tipo de dominio para uso en la capa de interfaces.
type MemberWithUsername = teamdomain.MemberWithUsername

// mapTeam convierte un EventTeam de dominio a su representacion JSON para la API.
func mapTeam(t *teamdomain.EventTeam) EventTeamResponse {
	return EventTeamResponse{
		ID:      t.ID.String(),
		EventID: t.EventID.String(),
		Name:    t.Name,
		Score:   t.Score,
	}
}
