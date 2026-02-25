package application

import (
	"errors"

	"github.com/google/uuid"

	"golabs-api/internal/eventteam/domain"
)

type LeaveTeamUseCase struct {
	teamsRepo domain.Repository
}

func NewLeaveTeamUseCase(
	teamsRepo domain.Repository,
) *LeaveTeamUseCase {
	return &LeaveTeamUseCase{teamsRepo: teamsRepo}
}

func (uc *LeaveTeamUseCase) Execute(
	teamID uuid.UUID,
	userID uuid.UUID,
) error {

	members, err := uc.teamsRepo.ListMembers(teamID)
	if err != nil {
		return err
	}

	var role domain.TeamRole
	for _, m := range members {
		if m.UserID == userID {
			role = m.Role
			break
		}
	}

	if role == "" {
		return errors.New("usuario no pertenece al equipo")
	}

	if role == domain.TeamOwner && len(members) > 1 {
		return errors.New("owner no puede salir si hay otros miembros")
	}

	return uc.teamsRepo.RemoveMember(teamID, userID)
}
