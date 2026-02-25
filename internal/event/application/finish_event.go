// Package application contiene los casos de uso del modulo de eventos.
package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

// FinishEventUseCase transiciona un evento del estado "running" a "finished".
// En estado "finished" no se aceptan nuevas flag submissions y el rankig es definitivo.
type FinishEventUseCase struct {
	repo eventdomain.Repository
}

// NewFinishEventUseCase crea un FinishEventUseCase con el repositorio indicado.
func NewFinishEventUseCase(repo eventdomain.Repository) *FinishEventUseCase {
	return &FinishEventUseCase{repo: repo}
}

// Execute aplica la transicion "running" -> "finished" al evento indicado.
// La validacion de la transicion de estado se hace en el metodo Event.Finish() del dominio.
//
// Entrada:  id, UUID del evento en formato string.
// Salida:   error si el UUID es invalido, el evento no existe o la transicion es invalida.
func (uc *FinishEventUseCase) Execute(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id invalido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	if err := event.Finish(); err != nil {
		return err
	}

	return uc.repo.Update(event)
}
