package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

type TeamRole string

const (
	TeamOwner  TeamRole = "owner"
	TeamMember TeamRole = "member"
)

type EventTeam struct {
	ID             uuid.UUID
	EventID        uuid.UUID
	Name           string
	JoinSecretHash string
	Score          int
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

type EventTeamMember struct {
	EventTeamID uuid.UUID
	UserID      uuid.UUID
	Role        TeamRole
	JoinedAt    time.Time
}

// Constructor
func NewEventTeam(
	eventID uuid.UUID,
	name string,
	joinSecretHash string,
) (*EventTeam, error) {

	if eventID == uuid.Nil {
		return nil, errors.New("eventID requerido")
	}

	if name == "" {
		return nil, errors.New("team name requerido")
	}

	if joinSecretHash == "" {
		return nil, errors.New("join secret inválido")
	}

	now := time.Now()

	return &EventTeam{
		ID:             uuid.New(),
		EventID:        eventID,
		Name:           name,
		JoinSecretHash: joinSecretHash,
		Score:          0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
