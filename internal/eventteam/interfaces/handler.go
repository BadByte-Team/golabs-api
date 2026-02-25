package interfaces

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"golabs-api/internal/apperrors"
	eventteamapp "golabs-api/internal/eventteam/application"
	authctx "golabs-api/internal/interfaces/http/middleware/auth"
)

type EventTeamHandler struct {
	createUC *eventteamapp.CreateTeamUseCase
	joinUC   *eventteamapp.JoinTeamUseCase
	leaveUC  *eventteamapp.LeaveTeamUseCase
	rotateUC *eventteamapp.RotateJoinSecretUseCase
}

func NewEventTeamHandler(
	create *eventteamapp.CreateTeamUseCase,
	join *eventteamapp.JoinTeamUseCase,
	leave *eventteamapp.LeaveTeamUseCase,
	rotate *eventteamapp.RotateJoinSecretUseCase,
) *EventTeamHandler {
	return &EventTeamHandler{
		createUC: create,
		joinUC:   join,
		leaveUC:  leave,
		rotateUC: rotate,
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	result, err := h.createUC.Execute(eventID, userID, req.Name)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusCreated, CreateTeamResponse{
		EventTeamResponse: EventTeamResponse{
			ID:    result.Team.ID.String(),
			Name:  result.Team.Name,
			Score: result.Team.Score,
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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
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
