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

type EventHandler struct {
	createUC *eventsapp.CreateEventUseCase
	getUC    *eventsapp.GetEventByIDUseCase
	listUC   *eventsapp.ListEventsUseCase
	openUC   *eventsapp.OpenEventUseCase
	startUC  *eventsapp.StartEventUseCase
	finishUC *eventsapp.FinishEventUseCase
}

func NewEventHandler(
	create *eventsapp.CreateEventUseCase,
	get *eventsapp.GetEventByIDUseCase,
	list *eventsapp.ListEventsUseCase,
	open *eventsapp.OpenEventUseCase,
	start *eventsapp.StartEventUseCase,
	finish *eventsapp.FinishEventUseCase,
) *EventHandler {
	return &EventHandler{
		createUC: create,
		getUC:    get,
		listUC:   list,
		openUC:   open,
		startUC:  start,
		finishUC: finish,
	}
}

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

func (h *EventHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")

	event, err := h.getUC.Execute(id)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, mapEvent(event))
}

func (h *EventHandler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)
	events, err := h.listUC.Execute()
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	// Apply manual pagination on the slice (list is small in practice).
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

func (h *EventHandler) Open(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	if err := h.openUC.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) Start(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	if err := h.startUC.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *EventHandler) Finish(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "event_id")
	if err := h.finishUC.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

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
