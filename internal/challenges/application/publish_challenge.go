// Package application contiene los casos de uso del modulo de challenges (retos CTF).
package application

import (
	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
)

// PublishChallengeUseCase cambia la visibilidad de un challenge (publicar o ocultar).
// Solo admins pueden ejecutar esta accion. Los challenges ocultos no son visibles para los equipos.
type PublishChallengeUseCase struct {
	repo challengedomain.Repository
}

// NewPublishChallengeUseCase crea un PublishChallengeUseCase con el repositorio indicado.
func NewPublishChallengeUseCase(repo challengedomain.Repository) *PublishChallengeUseCase {
	return &PublishChallengeUseCase{repo: repo}
}

// Execute cambia la visibilidad del challenge.
// Si publish=true, el challenge pasa a visible=true.
// Si publish=false, el challenge pasa a visible=false (ocultado).
//
// Entrada:  id (UUID del challenge), publish (true = publicar, false = ocultar).
// Salida:   puntero al Challenge actualizado o error.
func (uc *PublishChallengeUseCase) Execute(id uuid.UUID, publish bool) (*challengedomain.Challenge, error) {
	challenge, err := uc.repo.GetChallengeByID(id)
	if err != nil {
		return nil, err
	}

	if publish {
		challenge.Publish()
	} else {
		challenge.Unpublish()
	}

	if err := uc.repo.UpdateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
