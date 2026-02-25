package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type EventStatus string

const (
	EventDraft    EventStatus = "draft"
	EventOpen     EventStatus = "open"
	EventRunning  EventStatus = "running"
	EventFinished EventStatus = "finished"
)

type Event struct {
	ID          uuid.UUID
	Name        string
	Description string
	MaxTeamSize int
	Status      EventStatus
	StartsAt    time.Time
	EndsAt      time.Time
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Constructor (reglas mínimas)
func NewEvent(
	name string,
	description string,
	maxTeamSize int,
	startsAt, endsAt time.Time,
) (*Event, error) {

	if name == "" {
		return nil, errors.New("event name requerido")
	}

	if maxTeamSize <= 0 {
		return nil, errors.New("maxTeamSize inválido")
	}

	if endsAt.Before(startsAt) {
		return nil, errors.New("fechas inválidas")
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

// Reglas de negocio

func (e *Event) Open() error {
	if e.Status != EventDraft {
		return errors.New("solo eventos en draft pueden abrirse")
	}
	e.Status = EventOpen
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) Start() error {
	if e.Status != EventOpen {
		return errors.New("evento no está abierto")
	}
	e.Status = EventRunning
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) Finish() error {
	if e.Status != EventRunning {
		return errors.New("evento no está en curso")
	}
	e.Status = EventFinished
	e.UpdatedAt = time.Now()
	return nil
}

func (e *Event) IsOpen() bool {
	return e.Status == EventOpen
}
