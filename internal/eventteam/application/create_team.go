package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
	"golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

type CreateTeamUseCase struct {
	eventsRepo eventdomain.Repository
	teamsRepo  domain.Repository
}

type CreateTeamResult struct {
	Team       *domain.EventTeam
	JoinSecret string
}

func NewCreateTeamUseCase(
	eventsRepo eventdomain.Repository,
	teamsRepo domain.Repository,
) *CreateTeamUseCase {
	return &CreateTeamUseCase{
		eventsRepo: eventsRepo,
		teamsRepo:  teamsRepo,
	}
}

func (uc *CreateTeamUseCase) Execute(
	eventID uuid.UUID,
	ownerID uuid.UUID,
	teamName string,
) (*CreateTeamResult, error) {

	event, err := uc.eventsRepo.GetByID(eventID)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}

	if !event.IsOpen() {
		return nil, errors.New("evento no está abierto")
	}

	joinSecret, err := security.GenerateJoinSecret(10)
	if err != nil {
		return nil, err
	}

	hash := security.HashJoinSecret(joinSecret)

	team, err := domain.NewEventTeam(eventID, teamName, hash)
	if err != nil {
		return nil, err
	}

	if err := uc.teamsRepo.SaveTeam(team); err != nil {
		return nil, err
	}

	member := &domain.EventTeamMember{
		EventTeamID: team.ID,
		UserID:      ownerID,
		Role:        domain.TeamOwner,
	}

	if err := uc.teamsRepo.AddMember(member); err != nil {
		return nil, err
	}

	return &CreateTeamResult{
		Team:       team,
		JoinSecret: joinSecret,
	}, nil
}
