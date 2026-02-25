// Package application contiene los casos de uso del modulo de challenges (retos CTF).
package application

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"golabs-api/internal/apperrors"
	challengedomain "golabs-api/internal/challenges/domain"
	eventdomain "golabs-api/internal/event/domain"
	teamdomain "golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

// SubmitFlagResult indica el resultado de un intento de envio de flag.
type SubmitFlagResult struct {
	Correct bool // true si la flag era correcta
	Points  int  // puntos otorgados al equipo (0 si fue incorrecta o ya resuelta)
}

// SubmitFlagUseCase valida una flag enviada por un equipo y registra la resolucion si es correcta.
// Otorga los puntos del challenge al equipo en caso de acierto.
type SubmitFlagUseCase struct {
	challengeRepo challengedomain.Repository
	eventRepo     eventdomain.Repository
	teamRepo      teamdomain.Repository
}

// NewSubmitFlagUseCase crea un SubmitFlagUseCase con los repositorios indicados.
func NewSubmitFlagUseCase(
	challengeRepo challengedomain.Repository,
	eventRepo eventdomain.Repository,
	teamRepo teamdomain.Repository,
) *SubmitFlagUseCase {
	return &SubmitFlagUseCase{
		challengeRepo: challengeRepo,
		eventRepo:     eventRepo,
		teamRepo:      teamRepo,
	}
}

// Execute valida el intento de flag y registra el solve si es correcto.
//
// Flujo de validacion (en orden):
//  1. El evento debe existir y estar en estado "running".
//  2. El challenge debe existir y ser visible para participantes.
//  3. El usuario debe pertenecer a un equipo en el evento.
//  4. El equipo no debe haber resuelto ya el challenge (evita puntos dobles).
//  5. La flag hasheada debe coincidir con la almacenada.
//
// Seguridad: tanto "flag incorrecta" como "flag sin configurar" retornan
// Correct=false sin error, para no revelar el estado interno del challenge.
//
// Entrada:  challengeID, eventID, userID (UUIDs), attempt (flag en texto plano).
// Salida:   SubmitFlagResult con Correct y Points, o error de validacion/BD.
func (uc *SubmitFlagUseCase) Execute(
	challengeID uuid.UUID,
	eventID uuid.UUID,
	userID uuid.UUID,
	attempt string,
) (*SubmitFlagResult, error) {
	// 1. Verificar que el evento existe y esta en curso.
	event, err := uc.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("%w: evento no encontrado", apperrors.ErrNotFound)
	}
	if event.Status != eventdomain.EventRunning {
		return nil, errors.New("el evento no esta en curso")
	}

	// 2. Verificar que el challenge existe y esta visible para participantes.
	challenge, err := uc.challengeRepo.GetChallengeByID(challengeID)
	if err != nil {
		return nil, fmt.Errorf("%w: challenge no encontrado", apperrors.ErrNotFound)
	}
	if !challenge.Visible {
		return nil, fmt.Errorf("%w: challenge no encontrado", apperrors.ErrNotFound)
	}

	// 3. Encontrar el equipo del usuario en este evento.
	team, err := uc.teamRepo.GetTeamByUserAndEvent(eventID, userID)
	if err != nil {
		return nil, errors.New("debes pertenecer a un equipo para enviar una flag")
	}

	// 4. Verificar que el equipo no haya resuelto ya este challenge.
	alreadySolved, err := uc.challengeRepo.HasTeamSolved(challengeID, team.ID)
	if err != nil {
		return nil, err
	}
	if alreadySolved {
		return nil, fmt.Errorf("%w: tu equipo ya resolvio este challenge", apperrors.ErrConflict)
	}

	// 5. Comparar el hash del intento con el hash almacenado.
	// Si la flag no esta configurada o el hash no coincide, retornar incorrecto
	// con respuesta identica para no revelar el estado del challenge.
	flag, err := uc.challengeRepo.GetFlagByChallengeID(challengeID)
	if err != nil || security.Hash(attempt) != flag.Hash {
		return &SubmitFlagResult{Correct: false, Points: 0}, nil
	}

	// 6. Registrar el solve.
	solve := challengedomain.NewSolve(challengeID, team.ID, userID)
	if err := uc.challengeRepo.SaveSolve(solve); err != nil {
		return nil, err
	}

	// 7. Acumular los puntos del challenge al score del equipo.
	// Si el update falla, el solve ya fue guardado; se ignora el error
	// para no revertir el acierto del equipo.
	team.Score += challenge.Points
	team.UpdatedAt = solve.SolvedAt
	if err := uc.teamRepo.UpdateTeam(team); err != nil {
		_ = err // best-effort: el solve ya quedo registrado
	}

	return &SubmitFlagResult{Correct: true, Points: challenge.Points}, nil
}
