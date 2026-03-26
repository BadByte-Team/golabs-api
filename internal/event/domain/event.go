// Package domain define los tipos de dominio del modulo de eventos CTF,
// incluyendo el modelo Event, sus estados posibles y las reglas de negocio
// que controlan las transiciones de estado.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// EventStatus representa el Estado de ciclo de vida de un evento.
type EventStatus string

const (
	EventDraft    EventStatus = "draft"    // creado pero no abierto a inscripciones
	EventOpen     EventStatus = "open"     // aceptando equipos; aun no ha comenzado
	EventRunning  EventStatus = "running"  // en curso; se pueden enviar flags
	EventFinished EventStatus = "finished" // finalizado; no se aceptan mas submissions
)

// Event representa un evento CTF con sus fechas, estado y restricciones de equipo.
type Event struct {
	ID          uuid.UUID
	Name        string
	Description string
	MaxTeamSize int // numero maximo de miembros por equipo
	Status      EventStatus
	StartsAt    time.Time
	EndsAt      time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewEvent crea un Event en estado draft aplicando las reglas minimas de validacion.
//
// Reglas:
//   - name es obligatorio.
//   - maxTeamSize debe ser mayor que cero.
//   - endsAt debe ser posterior a startsAt.
//
// Retorna error si alguna regla se viola.
func NewEvent(
	name string,
	description string,
	maxTeamSize int,
	startsAt, endsAt time.Time,
) (*Event, error) {

	if name == "" {
		return nil, errors.New("el nombre del evento es requerido")
	}

	if maxTeamSize <= 0 {
		return nil, errors.New("maxTeamSize debe ser mayor que cero")
	}

	if endsAt.Before(startsAt) {
		return nil, errors.New("la fecha de fin debe ser posterior a la fecha de inicio")
	}

	now := time.Now()

	return &Event{
		ID:          uuid.New(),
		Name:        name,
		Description: description,
		MaxTeamSize: maxTeamSize,
		Status:      EventDraft,
		StartsAt:    startsAt,
		EndsAt:      endsAt,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Open transiciona el evento de draft a open para aceptar inscripciones de equipos.
// Retorna error si el evento no esta en estado draft.
func (e *Event) Open() error {
	if e.Status != EventDraft {
		return errors.New("solo eventos en estado draft pueden abrirse")
	}
	e.Status = EventOpen
	e.UpdatedAt = time.Now()
	return nil
}

// Start transiciona el evento de open a running para iniciar la competencia.
// Retorna error si el evento no esta en estado open.
func (e *Event) Start() error {
	if e.Status != EventOpen {
		return errors.New("el evento debe estar en estado open para iniciar")
	}
	e.Status = EventRunning
	e.UpdatedAt = time.Now()
	return nil
}

// Finish transiciona el evento de running a finished para cerrar la competencia.
// Retorna error si el evento no esta en estado running.
func (e *Event) Finish() error {
	if e.Status != EventRunning {
		return errors.New("el evento debe estar en curso para poder finalizarlo")
	}
	e.Status = EventFinished
	e.UpdatedAt = time.Now()
	return nil
}

// UpdateFields modifica los campos editables del evento.
// Solo se permite editar eventos en estado draft.
//
// Reglas:
//   - el evento debe estar en estado draft.
//   - name es obligatorio.
//   - maxTeamSize debe ser mayor que cero.
//   - endsAt debe ser posterior a startsAt.
//
// Retorna error si alguna regla se viola.
func (e *Event) UpdateFields(name, description string, maxTeamSize int, startsAt, endsAt time.Time) error {
	if e.Status != EventDraft {
		return errors.New("solo eventos en estado draft pueden editarse")
	}
	if name == "" {
		return errors.New("el nombre del evento es requerido")
	}
	if maxTeamSize <= 0 {
		return errors.New("maxTeamSize debe ser mayor que cero")
	}
	if endsAt.Before(startsAt) {
		return errors.New("la fecha de fin debe ser posterior a la fecha de inicio")
	}

	e.Name = name
	e.Description = description
	e.MaxTeamSize = maxTeamSize
	e.StartsAt = startsAt
	e.EndsAt = endsAt
	e.UpdatedAt = time.Now()
	return nil
}

// CanDelete retorna true si el evento puede ser eliminado (solo en estado draft).
func (e *Event) CanDelete() bool {
	return e.Status == EventDraft
}

// IsOpen retorna true si el evento acepta inscripciones de nuevos equipos.
func (e *Event) IsOpen() bool {
	return e.Status == EventOpen
}
