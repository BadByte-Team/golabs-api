// Package application contiene los casos de uso del modulo de eventos.
package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

// GetEventByIDUseCase obtiene un evento especifico por su UUID.
type GetEventByIDUseCase struct {
	repo eventdomain.Repository
}

// NewGetEventByIDUseCase crea un GetEventByIDUseCase con el repositorio indicado.
func NewGetEventByIDUseCase(repo eventdomain.Repository) *GetEventByIDUseCase {
	return &GetEventByIDUseCase{repo: repo}
}

// Execute busca y retorna el evento con el ID especificado.
//
// Entrada:  id, UUID del evento en formato string.
// Salida:   puntero al Event o error (apperrors.ErrNotFound si no existe).
func (uc *GetEventByIDUseCase) Execute(id string) (*eventdomain.Event, error) {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("id invalido")
	}

	return uc.repo.GetByID(eventID)
}
