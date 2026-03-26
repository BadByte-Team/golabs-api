// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de eventos.
package interfaces

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"golabs-api/internal/apperrors"
	eventsapp "golabs-api/internal/event/application"
	eventdomain "golabs-api/internal/event/domain"
	"golabs-api/internal/interfaces/http/pagination"
	"golabs-api/internal/interfaces/http/validate"
)

// EventHandler agrupa los handlers HTTP del modulo de eventos.
// Cada metodo corresponde a un endpoint de la API REST de eventos.
type EventHandler struct {
	createUC *eventsapp.CreateEventUseCase
	updateUC *eventsapp.UpdateEventUseCase
	deleteUC *eventsapp.DeleteEventUseCase
	getUC    *eventsapp.GetEventByIDUseCase
	listUC   *eventsapp.ListEventsUseCase
	openUC   *eventsapp.OpenEventUseCase
	startUC  *eventsapp.StartEventUseCase
	finishUC *eventsapp.FinishEventUseCase
}

// NewEventHandler inyecta las dependencias del EventHandler.
func NewEventHandler(
	create *eventsapp.CreateEventUseCase,
	update *eventsapp.UpdateEventUseCase,
	del *eventsapp.DeleteEventUseCase,
	get *eventsapp.GetEventByIDUseCase,
	list *eventsapp.ListEventsUseCase,
	open *eventsapp.OpenEventUseCase,
	start *eventsapp.StartEventUseCase,
	finish *eventsapp.FinishEventUseCase,
) *EventHandler {
	return &EventHandler{
		createUC: create,
		updateUC: update,
		deleteUC: del,
		getUC:    get,
		listUC:   list,
		openUC:   open,
		startUC:  start,
		finishUC: finish,
	}
}

// Create godoc
//
// POST /api/v1/events
//
// Crea un nuevo evento en estado "draft". Solo admins.
// Body:  CreateEventRequest
// Exito: 201 EventResponse
func (h *EventHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req CreateEventRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	event, err := h.createUC.Execute(
		req.Name,
		req.Description,
		req.MaxTeamSize,
		req.StartsAt,
		req.EndsAt,
	)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusCreated, mapEvent(event))
}

// Update godoc — PUT /api/v1/events/{event_id}
// Actualiza los datos de un evento existente. Solo admins. Solo eventos en estado draft.
// Body: UpdateEventRequest | Exito: 200 EventResponse
func (h *EventHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")

	var req UpdateEventRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondJSON(w, http.StatusUnprocessableEntity, map[string]string{"error": err.Error()})
		return
	}

	event, err := h.updateUC.Execute(
		id,
		req.Name,
		req.Description,
		req.MaxTeamSize,
		req.StartsAt,
		req.EndsAt,
	)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, mapEvent(event))
}

// Delete godoc — POST /api/v1/events/{event_id}/delete
// Elimina un evento existente. Solo admins. Solo eventos en estado draft.
// Exito: 204 No Content
func (h *EventHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	if err := h.deleteUC.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetByID godoc
//
// GET /api/v1/events/{event_id}
//
// Retorna el detalle de un evento especifico.
// Exito: 200 EventResponse
// Error: 404 si no existe
func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")

	event, err := h.getUC.Execute(id)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, mapEvent(event))
}

// List godoc
//
// GET /api/v1/events?page=1&size=20
//
// Retorna todos los eventos paginados. La paginacion se aplica en memoria sobre el slice completo.
// Exito: 200 pagination.Response[EventResponse]
func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)
	events, err := h.listUC.Execute()
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	// Paginacion manual sobre el slice (lista reducida en la practica).
	total := len(events)
	start := pg.Offset()
	if start > total {
		start = total
	}
	end := start + pg.Size
	if end > total {
		end = total
	}
	page := events[start:end]

	resp := make([]EventResponse, 0, len(page))
	for _, e := range page {
		resp = append(resp, mapEvent(e))
	}

	apperrors.RespondJSON(w, http.StatusOK, pagination.New(resp, pg, total))
}

// Open godoc — PATCH /api/v1/events/{event_id}/open
// Transiciona el evento de "draft" a "open" para que los equipos puedan unirse.
// Exito: 204 No Content
func (h *EventHandler) Open(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	if err := h.openUC.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Start godoc — PATCH /api/v1/events/{event_id}/start
// Transiciona el evento de "open" a "running". Los challenges quedan visibles para los equipos.
// Exito: 204 No Content
func (h *EventHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	if err := h.startUC.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// Finish godoc — PATCH /api/v1/events/{event_id}/finish
// Transiciona el evento de "running" a "finished". No se aceptan mas flag submissions.
// Exito: 204 No Content
func (h *EventHandler) Finish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	if err := h.finishUC.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// mapEvent convierte un Event de dominio a su representacion JSON para la API.
func mapEvent(e *eventdomain.Event) EventResponse {
	return EventResponse{
		ID:          e.ID.String(),
		Name:        e.Name,
		Description: e.Description,
		MaxTeamSize: e.MaxTeamSize,
		Status:      string(e.Status),
		StartsAt:    e.StartsAt,
		EndsAt:      e.EndsAt,
		CreatedAt:   e.CreatedAt,
		UpdatedAt:   e.UpdatedAt,
	}
}
