// Package application contiene los casos de uso del modulo de challenges (retos CTF).
package application

import (
	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
)

// UpdateChallengeUseCase modifica los datos editables de un challenge existente.
// Un challenge puede ser actualizado en cualquier estado (visible o no).
type UpdateChallengeUseCase struct {
	repo challengedomain.Repository
}

// NewUpdateChallengeUseCase crea un UpdateChallengeUseCase con el repositorio indicado.
func NewUpdateChallengeUseCase(repo challengedomain.Repository) *UpdateChallengeUseCase {
	return &UpdateChallengeUseCase{repo: repo}
}

// Execute actualiza los campos modificables del challenge.
// La validacion de los campos (titulo no vacio, puntos >= 0, etc.) se
// delega al metodo Challenge.Update() del dominio para mantener las reglas de negocio centralizadas.
//
// Entrada:  id (UUID del challenge), title, description, category, points, difficulty.
// Salida:   puntero al Challenge actualizado o error de validacion/BD.
func (uc *UpdateChallengeUseCase) Execute(
	id uuid.UUID,
	title, description string,
	category challengedomain.ChallengeCategory,
	points int,
	difficulty challengedomain.ChallengeDifficulty,
	fileURL string,
) (*challengedomain.Challenge, error) {
	challenge, err := uc.repo.GetChallengeByID(id)
	if err != nil {
		return nil, err
	}

	if err := challenge.Update(title, description, category, points, difficulty, fileURL); err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
