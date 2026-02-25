// Package application contiene los casos de uso del modulo de challenges (retos CTF).
package application

import (
	"time"

	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
	"golabs-api/internal/infrastructure/security"
)

// SetFlagUseCase establece o reemplaza la flag de un challenge.
// La flag se almacena SIEMPRE como hash SHA-256; el texto plano nunca se persiste.
// Si ya existia una flag para el challenge, se reemplaza (upsert).
type SetFlagUseCase struct {
	repo challengedomain.Repository
}

// NewSetFlagUseCase crea un SetFlagUseCase con el repositorio indicado.
func NewSetFlagUseCase(repo challengedomain.Repository) *SetFlagUseCase {
	return &SetFlagUseCase{repo: repo}
}

// Execute hashea la flag en texto plano y la persiste mediante upsert.
//
// Solo el hash se almacena en BD para que si la base de datos es comprometida,
// las flags de los challenges no puedan ser recuperadas.
//
// Entrada:  challengeID (UUID del challenge), plaintext (flag en texto plano).
// Salida:   error si el challenge no existe o falla el upsert en BD.
func (uc *SetFlagUseCase) Execute(challengeID uuid.UUID, plaintext string) error {
	// Verificar que el challenge existe antes de crear la flag.
	if _, err := uc.repo.GetChallengeByID(challengeID); err != nil {
		return err
	}

	flag := &challengedomain.Flag{
		ID:          uuid.New(),
		ChallengeID: challengeID,
		Hash:        security.Hash(plaintext), // solo el hash va a BD
		CreatedAt:   time.Now(),
	}

	return uc.repo.UpsertFlag(flag)
}
