// Package application contiene los casos de uso del modulo de eventos.
package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

// StartEventUseCase transiciona un evento del estado "open" a "running".
// En estado "running", los challenges son visibles para los equipos y se aceptan flag submissions.
type StartEventUseCase struct {
	repo eventdomain.Repository
}

// NewStartEventUseCase crea un StartEventUseCase con el repositorio indicado.
func NewStartEventUseCase(repo eventdomain.Repository) *StartEventUseCase {
	return &StartEventUseCase{repo: repo}
}

// Execute aplica la transicion "open" -> "running" al evento indicado.
// La validacion de la transicion de estado se hace en el metodo Event.Start() del dominio.
//
// Entrada:  id, UUID del evento en formato string.
// Salida:   error si el UUID es invalido, el evento no existe o la transicion es invalida.
func (uc *StartEventUseCase) Execute(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id invalido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	if err := event.Start(); err != nil {
		return err
	}

	return uc.repo.Update(event)
}
