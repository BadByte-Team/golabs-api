// Package application contiene los casos de uso del modulo de eventos.
package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

// OpenEventUseCase transiciona un evento del estado "draft" a "open".
// En estado "open", los equipos pueden inscribirse al evento.
type OpenEventUseCase struct {
	repo eventdomain.Repository
}

// NewOpenEventUseCase crea un OpenEventUseCase con el repositorio indicado.
func NewOpenEventUseCase(repo eventdomain.Repository) *OpenEventUseCase {
	return &OpenEventUseCase{repo: repo}
}

// Execute aplica la transicion "draft" -> "open" al evento indicado.
// La validacion de la transicion de estado se hace en el metodo Event.Open() del dominio.
//
// Entrada:  id, UUID del evento en formato string.
// Salida:   error si el UUID es invalido, el evento no existe o la transicion es invalida.
func (uc *OpenEventUseCase) Execute(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id invalido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	if err := event.Open(); err != nil {
		return err
	}

	return uc.repo.Update(event)
}
