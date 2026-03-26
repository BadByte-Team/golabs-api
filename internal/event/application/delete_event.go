// Package application contiene los casos de uso del modulo de eventos.
package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

// DeleteEventUseCase elimina un evento existente.
// Solo se permite eliminar eventos en estado "draft".
type DeleteEventUseCase struct {
	repo eventdomain.Repository
}

// NewDeleteEventUseCase crea un DeleteEventUseCase con el repositorio indicado.
func NewDeleteEventUseCase(repo eventdomain.Repository) *DeleteEventUseCase {
	return &DeleteEventUseCase{repo: repo}
}

// Execute elimina el evento indicado si esta en estado draft.
//
// Entrada:  id, UUID del evento en formato string.
// Salida:   error si el UUID es invalido, el evento no existe o no esta en estado draft.
func (uc *DeleteEventUseCase) Execute(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id invalido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	if !event.CanDelete() {
		return errors.New("solo eventos en estado draft pueden eliminarse")
	}

	return uc.repo.Delete(eventID)
}
