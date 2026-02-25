package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
	"golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

type JoinTeamUseCase struct {
	eventsRepo eventdomain.Repository
	teamsRepo  domain.Repository
}

func NewJoinTeamUseCase(
	eventsRepo eventdomain.Repository,
	teamsRepo domain.Repository,
) *JoinTeamUseCase {
	return &JoinTeamUseCase{
		eventsRepo: eventsRepo,
		teamsRepo:  teamsRepo,
	}
}

func (uc *JoinTeamUseCase) Execute(
	eventID uuid.UUID,
	userID uuid.UUID,
	teamName string,
	joinSecret string,
) error {

	event, err := uc.eventsRepo.GetByID(eventID)
	if err != nil {
		return errors.New("evento no encontrado")
	}

	if !event.IsOpen() {
		return errors.New("evento no está abierto")
	}

	inEvent, err := uc.teamsRepo.IsUserInEvent(eventID, userID)
	if err != nil {
		return err
	}
	if inEvent {
		return errors.New("usuario ya pertenece a un equipo en este evento")
	}

	team, err := uc.teamsRepo.GetTeamByName(eventID, teamName)
	if err != nil {
		return errors.New("no se pudo unir al equipo")
	}

	if security.HashJoinSecret(joinSecret) != team.JoinSecretHash {
		return errors.New("no se pudo unir al equipo")
	}

	count, err := uc.teamsRepo.CountMembers(team.ID)
	if err != nil {
		return err
	}
	if count >= event.MaxTeamSize {
		return errors.New("equipo lleno")
	}

	member := &domain.EventTeamMember{
		EventTeamID: team.ID,
		UserID:      userID,
		Role:        domain.TeamMember,
	}

	return uc.teamsRepo.AddMember(member)
}
