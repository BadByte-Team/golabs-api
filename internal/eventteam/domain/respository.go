package domain

import (
	"time"

	"github.com/google/uuid"
)

// MemberWithUsername enriches EventTeamMember with the user's display name,
// resolved at query time to avoid extra round-trips in the handler.
type MemberWithUsername struct {
	EventTeamID uuid.UUID
	UserID      uuid.UUID
	Username    string
	Role        TeamRole
	JoinedAt    time.Time
}

// LeaderboardEntry is the read-model for the public leaderboard.
type LeaderboardEntry struct {
	Rank        int    `json:"rank"`
	TeamID      string `json:"team_id"`
	TeamName    string `json:"team_name"`
	Score       int    `json:"score"`
	MemberCount int    `json:"member_count"`
}

type Repository interface {
	SaveTeam(team *EventTeam) error
	UpdateTeam(team *EventTeam) error
	GetTeamByID(id uuid.UUID) (*EventTeam, error)
	GetTeamByName(eventID uuid.UUID, name string) (*EventTeam, error)
	// GetTeamByUserAndEvent returns the team the given user belongs to in an event.
	GetTeamByUserAndEvent(eventID, userID uuid.UUID) (*EventTeam, error)

	// ListTeamsByEvent returns all teams in an event sorted by score DESC.
	ListTeamsByEvent(eventID uuid.UUID) ([]*EventTeam, error)

	AddMember(member *EventTeamMember) error
	RemoveMember(eventTeamID, userID uuid.UUID) error
	ListMembers(eventTeamID uuid.UUID) ([]*EventTeamMember, error)
	// ListMembersWithUsername joins members with the users table for display.
	ListMembersWithUsername(eventTeamID uuid.UUID) ([]*MemberWithUsername, error)

	CountMembers(eventTeamID uuid.UUID) (int, error)
	IsUserInEvent(eventID, userID uuid.UUID) (bool, error)
}
