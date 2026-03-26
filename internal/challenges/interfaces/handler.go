// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de challenges.
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

// ChallengeHandler agrupa los handlers HTTP del modulo de challenges.
type ChallengeHandler struct {
	createUC  *challengeapp.CreateChallengeUseCase
	updateUC  *challengeapp.UpdateChallengeUseCase
	publishUC *challengeapp.PublishChallengeUseCase
	listUC    *challengeapp.ListChallengesUseCase
	getUC     *challengeapp.GetChallengeUseCase
	setFlagUC *challengeapp.SetFlagUseCase
	submitUC  *challengeapp.SubmitFlagUseCase
}

// NewChallengeHandler inyecta las dependencias del ChallengeHandler.
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

// List godoc — GET /api/v1/events/{event_id}/challenges?category=&difficulty=
// Lista los challenges del evento. Admins ven todo (incluidos ocultos); participantes solo los visibles.
// Exito: 200 []ChallengeResponse (con SolveCount y FirstBloodTeamID)
func (h *ChallengeHandler) List(w http.ResponseWriter, r *http.Request) {
	eventID, err := uuid.Parse(chi.URLParam(r, "event_id"))
	if err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	user, _ := authmw.GetUser(r.Context())
	isAdmin := user.Role == userdomain.RoleAdmin

	// Filtros opcionales por categoria y dificultad.
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

// Get godoc — GET /api/v1/events/{event_id}/challenges/{challenge_id}
// Retorna el detalle de un challenge. Los ocultos retornan 404 para no-admins.
// Exito: 200 ChallengeResponse
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

// Create godoc — POST /api/v1/events/{event_id}/challenges (admin)
// Crea un challenge en estado oculto (visible=false). Requiere SetFlag para activarlo.
// Body: CreateChallengeRequest | Exito: 201 ChallengeResponse
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
		req.FileURL,
	)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusCreated, mapChallenge(challenge))
}

// Update godoc — PUT /api/v1/events/{event_id}/challenges/{challenge_id} (admin)
// Actualiza los datos de un challenge. La flag se gestiona por separado con SetFlag.
// Body: UpdateChallengeRequest | Exito: 200 ChallengeResponse
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
		req.FileURL,
	)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapChallenge(challenge))
}

// Publish godoc — POST /api/v1/events/{event_id}/challenges/{challenge_id}/publish (admin)
// Hace visible el challenge para los participantes.
// Exito: 200 ChallengeResponse
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

// Unpublish godoc — POST /api/v1/events/{event_id}/challenges/{challenge_id}/unpublish (admin)
// Oculta el challenge para los participantes (visible=false).
// Exito: 200 ChallengeResponse
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

// SetFlag godoc — POST /api/v1/events/{event_id}/challenges/{challenge_id}/flag (admin)
// Establece o reemplaza la flag del challenge. El texto plano se hashea en el servidor.
// Body: SetFlagRequest | Exito: 204 No Content
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

// Submit godoc — POST /api/v1/events/{event_id}/challenges/{challenge_id}/submit
// Valida una flag enviada por el equipo. La respuesta es intencionalmente vaga para
// no revelar informacion sobre la flag correcta.
// Body: SubmitFlagRequest | Exito: 200 SubmitFlagResponse
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

// mapChallenge convierte un Challenge de dominio a su representacion JSON para la API.
func mapChallenge(c *challengedomain.Challenge) ChallengeResponse {
	return ChallengeResponse{
		ID:          c.ID.String(),
		EventID:     c.EventID.String(),
		Title:       c.Title,
		Description: c.Description,
		Category:    string(c.Category),
		Points:      c.Points,
		Difficulty:  string(c.Difficulty),
		FileURL:     c.FileURL,
		Visible:     c.Visible,
		CreatedAt:   c.CreatedAt,
		UpdatedAt:   c.UpdatedAt,
	}
}

// mapChallengeResult enriquece la respuesta del challenge con estadisticas de resoluciones.
func mapChallengeResult(res *challengeapp.ListChallengesResult) ChallengeResponse {
	r := mapChallenge(res.Challenge)
	r.SolveCount = res.SolveCount
	if res.FirstBlood != nil {
		tid := res.FirstBlood.EventTeamID.String()
		r.FirstBloodTeamID = &tid
	}
	return r
}
