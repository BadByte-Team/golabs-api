// Package application contiene los casos de uso del modulo de challenges (retos CTF).
package application

import (
	"fmt"

	"github.com/google/uuid"

	"golabs-api/internal/apperrors"
	challengedomain "golabs-api/internal/challenges/domain"
)

// GetChallengeUseCase obtiene un challenge especifico por su UUID.
// Los challenges ocultos (visible=false) son invisibles para usuarios no admin.
type GetChallengeUseCase struct {
	repo challengedomain.Repository
}

// NewGetChallengeUseCase crea un GetChallengeUseCase con el repositorio indicado.
func NewGetChallengeUseCase(repo challengedomain.Repository) *GetChallengeUseCase {
	return &GetChallengeUseCase{repo: repo}
}

// Execute retorna el challenge con el ID indicado.
// Si el challenge no existe o esta oculto y el solicitante no es admin, retorna ErrNotFound.
// Esto evita que usuarios no admin puedan inferir la existencia de challenges ocultos.
//
// Entrada:  id (UUID del challenge), isAdmin (bool que indica si el solicitante es admin).
// Salida:   puntero al Challenge o error (ErrNotFound en ambos casos de fallo para no-admins).
func (uc *GetChallengeUseCase) Execute(id uuid.UUID, isAdmin bool) (*challengedomain.Challenge, error) {
	challenge, err := uc.repo.GetChallengeByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w", apperrors.ErrNotFound)
	}

	// Tratar challenges ocultos como inexistentes para usuarios no admin.
	if !isAdmin && !challenge.Visible {
		return nil, fmt.Errorf("%w: challenge not found", apperrors.ErrNotFound)
	}

	return challenge, nil
}
