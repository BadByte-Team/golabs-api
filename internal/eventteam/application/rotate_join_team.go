package application

import (
	"errors"

	"github.com/google/uuid"

	"golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

type RotateJoinSecretUseCase struct {
	teamsRepo domain.Repository
}

func NewRotateJoinSecretUseCase(
	teamsRepo domain.Repository,
) *RotateJoinSecretUseCase {
	return &RotateJoinSecretUseCase{teamsRepo: teamsRepo}
}

func (uc *RotateJoinSecretUseCase) Execute(
	teamID uuid.UUID,
	requesterID uuid.UUID,
) (string, error) {

	members, err := uc.teamsRepo.ListMembers(teamID)
	if err != nil {
		return "", err
	}

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
