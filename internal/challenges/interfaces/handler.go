package interfaces

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/google/uuid"

	"golabs-api/internal/apperrors"
	challengeapp "golabs-api/internal/challenges/application"
	challengedomain "golabs-api/internal/challenges/domain"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/validate"
	userdomain "golabs-api/internal/user/domain"
)

type ChallengeHandler struct {
	createUC  *challengeapp.CreateChallengeUseCase
	updateUC  *challengeapp.UpdateChallengeUseCase
	publishUC *challengeapp.PublishChallengeUseCase
	listUC    *challengeapp.ListChallengesUseCase
	getUC     *challengeapp.GetChallengeUseCase
	setFlagUC *challengeapp.SetFlagUseCase
	submitUC  *challengeapp.SubmitFlagUseCase
}

func NewChallengeHandler(
	create *challengeapp.CreateChallengeUseCase,
	update *challengeapp.UpdateChallengeUseCase,
	publish *challengeapp.PublishChallengeUseCase,
	list *challengeapp.ListChallengesUseCase,
	get *challengeapp.GetChallengeUseCase,
	setFlag *challengeapp.SetFlagUseCase,
	submit *challengeapp.SubmitFlagUseCase,
) *ChallengeHandler {
	return &ChallengeHandler{
		createUC:  create,
		updateUC:  update,
		publishUC: publish,
		listUC:    list,
		getUC:     get,
		setFlagUC: setFlag,
		submitUC:  submit,
	}
}

/* ── List ───────────────────────────────────────────────────────────────── */

func (h *ChallengeHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	user, _ := authmw.GetUser(r.Context())
	isAdmin := user.Role == userdomain.RoleAdmin

	// Optional ?category= and ?difficulty= filters
	category := r.URL.Query().Get("category")
	difficulty := r.URL.Query().Get("difficulty")

	results, err := h.listUC.Execute(eventID, isAdmin, category, difficulty)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	resp := make([]ChallengeResponse, 0, len(results))
	for _, res := range results {
		resp = append(resp, mapChallengeResult(res))
	}
	apperrors.RespondJSON(w, http.StatusOK, resp)
}

/* ── Get ────────────────────────────────────────────────────────────────── */

func (h *ChallengeHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "challenge_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	user, _ := authmw.GetUser(r.Context())
	isAdmin := user.Role == userdomain.RoleAdmin

	challenge, err := h.getUC.Execute(id, isAdmin)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapChallenge(challenge))
}

/* ── Create (admin) ─────────────────────────────────────────────────────── */

func (h *ChallengeHandler) Create(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	var req CreateChallengeRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	challenge, err := h.createUC.Execute(
		eventID,
		req.Title,
		req.Description,
		challengedomain.ChallengeCategory(req.Category),
		req.Points,
		challengedomain.ChallengeDifficulty(req.Difficulty),
	)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusCreated, mapChallenge(challenge))
}

/* ── Update (admin) ─────────────────────────────────────────────────────── */

func (h *ChallengeHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "challenge_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	var req UpdateChallengeRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	challenge, err := h.updateUC.Execute(
		id,
		req.Title,
		req.Description,
		challengedomain.ChallengeCategory(req.Category),
		req.Points,
		challengedomain.ChallengeDifficulty(req.Difficulty),
	)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapChallenge(challenge))
}

/* ── Publish (admin) ─────────────────────────────────────────────────────── */

func (h *ChallengeHandler) Publish(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "challenge_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	challenge, err := h.publishUC.Execute(id, true)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapChallenge(challenge))
}

func (h *ChallengeHandler) Unpublish(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "challenge_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	challenge, err := h.publishUC.Execute(id, false)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapChallenge(challenge))
}

/* ── SetFlag (admin) ─────────────────────────────────────────────────────── */

func (h *ChallengeHandler) SetFlag(w http.ResponseWriter, r *http.Request) {
	id, err := uuid.Parse(chi.URLParam(r, "challenge_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	var req SetFlagRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	if err := h.setFlagUC.Execute(id, req.Flag); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

/* ── Submit (participant) ────────────────────────────────────────────────── */

func (h *ChallengeHandler) Submit(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	challengeID, err := uuid.Parse(chi.URLParam(r, "challenge_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	user, _ := authmw.GetUser(r.Context())
	userID, err := uuid.Parse(user.UserID)
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrUnauthorized)
		return
	}

	var req SubmitFlagRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	result, err := h.submitUC.Execute(challengeID, eventID, userID, req.Flag)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, SubmitFlagResponse{
		Correct: result.Correct,
		Points:  result.Points,
	})
}

/* ── helpers ─────────────────────────────────────────────────────────────── */

func mapChallenge(c *challengedomain.Challenge) ChallengeResponse {
	return ChallengeResponse{
		ID:          c.ID.String(),
		EventID:     c.EventID.String(),
		Title:       c.Title,
		Description: c.Description,
		Category:    string(c.Category),
		Points:      c.Points,
		Difficulty:  string(c.Difficulty),
		Visible:     c.Visible,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// mapChallengeResult enriches the response with solve stats.
func mapChallengeResult(res *challengeapp.ListChallengesResult) ChallengeResponse {
	r := mapChallenge(res.Challenge)
	r.SolveCount = res.SolveCount
	if res.FirstBlood != nil {
		tid := res.FirstBlood.EventTeamID.String()
		r.FirstBloodTeamID = &tid
	}
	return r
}
