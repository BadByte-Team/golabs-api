package application

import (
	"github.com/google/uuid"

	teamdomain "golabs-api/internal/eventteam/domain"
)

// ListTeamsByEventUseCase returns all teams in an event.
type ListTeamsByEventUseCase struct {
	repo teamdomain.Repository
}

func NewListTeamsByEventUseCase(repo teamdomain.Repository) *ListTeamsByEventUseCase {
	return &ListTeamsByEventUseCase{repo: repo}
}

func (uc *ListTeamsByEventUseCase) Execute(eventID uuid.UUID) ([]*teamdomain.EventTeam, error) {
	return uc.repo.ListTeamsByEvent(eventID)
}

// ExecuteMembers returns all members of a team with their usernames resolved.
func (uc *ListTeamsByEventUseCase) ExecuteMembers(teamID uuid.UUID) ([]*teamdomain.MemberWithUsername, error) {
	return uc.repo.ListMembersWithUsername(teamID)
}
