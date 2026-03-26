// Package application contiene los casos de uso del modulo de eventos.
package application

import (
	"errors"
	"time"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

// UpdateEventUseCase modifica los datos editables de un evento existente.
// Solo se permite editar eventos en estado "draft".
type UpdateEventUseCase struct {
	repo eventdomain.Repository
}

// NewUpdateEventUseCase crea un UpdateEventUseCase con el repositorio indicado.
func NewUpdateEventUseCase(repo eventdomain.Repository) *UpdateEventUseCase {
	return &UpdateEventUseCase{repo: repo}
}

// Execute actualiza los campos modificables del evento.
// La validacion de los campos y del estado se delega al metodo Event.UpdateFields() del dominio.
//
// Entrada:  id (UUID string), name, description, maxTeamSize, startsAt, endsAt.
// Salida:   puntero al Event actualizado o error de validacion/BD.
func (uc *UpdateEventUseCase) Execute(
	id string,
	name string,
	description string,
	maxTeamSize int,
	startsAt time.Time,
	endsAt time.Time,
) (*eventdomain.Event, error) {

	eventID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("id invalido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return nil, err
	}

	if err := event.UpdateFields(name, description, maxTeamSize, startsAt, endsAt); err != nil {
		return nil, err
	}

	if err := uc.repo.Update(event); err != nil {
		return nil, err
	}

	return event, nil
}
