package domain

import "github.com/google/uuid"

type Repository interface {
	SaveTeam(team *EventTeam) error
	UpdateTeam(team *EventTeam) error
	GetTeamByID(id uuid.UUID) (*EventTeam, error)
	GetTeamByName(eventID uuid.UUID, name string) (*EventTeam, error)
	// GetTeamByUserAndEvent returns the team the given user belongs to in an event.
	GetTeamByUserAndEvent(eventID, userID uuid.UUID) (*EventTeam, error)

	AddMember(member *EventTeamMember) error
	RemoveMember(eventTeamID, userID uuid.UUID) error
	ListMembers(eventTeamID uuid.UUID) ([]*EventTeamMember, error)

	CountMembers(eventTeamID uuid.UUID) (int, error)
	IsUserInEvent(eventID, userID uuid.UUID) (bool, error)
}
