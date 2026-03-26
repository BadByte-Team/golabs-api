// Package application contiene los casos de uso del modulo de challenges (retos CTF).
package application

import (
	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
	eventdomain "golabs-api/internal/event/domain"
)

// CreateChallengeUseCase crea un nuevo challenge en el sistema.
// El challenge se crea en estado "oculto" (visible=false) hasta que el admin lo publique.
// La flag no se establece en este paso; se configura con SetFlagUseCase.
type CreateChallengeUseCase struct {
	challengeRepo challengedomain.Repository
	eventRepo     eventdomain.Repository
}

// NewCreateChallengeUseCase crea un CreateChallengeUseCase con los repositorios indicados.
func NewCreateChallengeUseCase(
	challengeRepo challengedomain.Repository,
	eventRepo eventdomain.Repository,
) *CreateChallengeUseCase {
	return &CreateChallengeUseCase{
		challengeRepo: challengeRepo,
		eventRepo:     eventRepo,
	}
}

// Execute crea y persiste un nuevo challenge para el evento indicado.
//
// Verifica que el evento exista antes de crear el challenge.
// La validacion del titulo, descripcion y puntos se delega al constructor de dominio.
//
// Entrada:  eventID, title, description, category, points, difficulty.
// Salida:   puntero al Challenge creado (visible=false) o error de validacion/BD.
func (uc *CreateChallengeUseCase) Execute(
	eventID uuid.UUID,
	title, description string,
	category challengedomain.ChallengeCategory,
	points int,
	difficulty challengedomain.ChallengeDifficulty,
	fileURL string,
) (*challengedomain.Challenge, error) {
	// Verificar que el evento existe antes de asociar el challenge.
	if _, err := uc.eventRepo.GetByID(eventID); err != nil {
		return nil, err
	}

	challenge, err := challengedomain.NewChallenge(eventID, title, description, category, points, difficulty, fileURL)
	if err != nil {
		return nil, err
	}

	if err := uc.challengeRepo.SaveChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
