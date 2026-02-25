// Package application contiene los casos de uso del modulo de equipos por evento.
package application

import (
	"errors"

	"github.com/google/uuid"

	"golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

// RotateJoinSecretUseCase genera un nuevo join secret para el equipo e invalida el anterior.
// Solo el owner del equipo puede rotar el secreto.
//
// Util cuando se sospecha que el join secret fue comprometido (por ejemplo, si
// alguien comparte el link de invitacion con personas no deseadas).
type RotateJoinSecretUseCase struct {
	teamsRepo domain.Repository
}

// NewRotateJoinSecretUseCase crea un RotateJoinSecretUseCase con el repositorio indicado.
func NewRotateJoinSecretUseCase(
	teamsRepo domain.Repository,
) *RotateJoinSecretUseCase {
	return &RotateJoinSecretUseCase{teamsRepo: teamsRepo}
}

// Execute verifica que el solicitante sea el owner y genera un nuevo join secret.
//
// El join secret anterior queda invalidado inmediatamente al persistir el nuevo hash.
//
// Entrada:  teamID, requesterID (UUIDs del equipo y el solicitante).
// Salida:   nuevo join secret en texto plano (para mostrarlo al owner UNA SOLA VEZ), o error.
func (uc *RotateJoinSecretUseCase) Execute(
	teamID uuid.UUID,
	requesterID uuid.UUID,
) (string, error) {

	members, err := uc.teamsRepo.ListMembers(teamID)
	if err != nil {
		return "", err
	}

	// Verificar que el solicitante es el owner del equipo.
	isOwner := false
	for _, m := range members {
		if m.UserID == requesterID && m.Role == domain.TeamOwner {
			isOwner = true
			break
		}
	}
	if !isOwner {
		return "", errors.New("solo el owner puede rotar el secreto")
	}

	secret, err := security.GenerateJoinSecret(10)
	if err != nil {
		return "", err
	}

	hash := security.HashJoinSecret(secret)

	team, err := uc.teamsRepo.GetTeamByID(teamID)
	if err != nil {
		return "", err
	}

	team.JoinSecretHash = hash
	return secret, uc.teamsRepo.UpdateTeam(team)
}
