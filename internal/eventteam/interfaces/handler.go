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

type EventTeamHandler struct {
	createUC      *eventteamapp.CreateTeamUseCase
	joinUC        *eventteamapp.JoinTeamUseCase
	leaveUC       *eventteamapp.LeaveTeamUseCase
	rotateUC      *eventteamapp.RotateJoinSecretUseCase
	listTeamsUC   *eventteamapp.ListTeamsByEventUseCase
	leaderboardUC *eventteamapp.GetLeaderboardUseCase
}

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

// ListTeams handles GET /events/{event_id}/teams
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

// ListMembers handles GET /events/{event_id}/teams/{team_id}/members
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

// Leaderboard handles GET /events/{event_id}/leaderboard
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

/* ── helpers ─────────────────────────────────────────────────────────────── */

// MemberWithUsername extends EventTeamMember with resolved username.
// Used internally by ListTeamsByEventUseCase.ExecuteMembers.
type MemberWithUsername = teamdomain.MemberWithUsername

func mapTeam(t *teamdomain.EventTeam) EventTeamResponse {
	return EventTeamResponse{
		ID:      t.ID.String(),
		EventID: t.EventID.String(),
		Name:    t.Name,
		Score:   t.Score,
	}
}
