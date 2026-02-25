package domain

import "github.com/google/uuid"

type Repository interface {
	// ── Challenges ──────────────────────────────────────────────────────────

	SaveChallenge(c *Challenge) error
	UpdateChallenge(c *Challenge) error
	GetChallengeByID(id uuid.UUID) (*Challenge, error)
	// visibleOnly=true for participants; false for admins
	ListChallengesByEvent(eventID uuid.UUID, visibleOnly bool) ([]*Challenge, error)

	// ── Flags ───────────────────────────────────────────────────────────────

	// UpsertFlag inserts or replaces the flag for a challenge (1 flag per challenge).
	UpsertFlag(f *Flag) error
	GetFlagByChallengeID(challengeID uuid.UUID) (*Flag, error)

	// ── Solves ──────────────────────────────────────────────────────────────

	SaveSolve(s *Solve) error
	HasTeamSolved(challengeID, teamID uuid.UUID) (bool, error)
	ListSolvesByChallenge(challengeID uuid.UUID) ([]*Solve, error)
	ListSolvesByTeam(teamID uuid.UUID) ([]*Solve, error)
}
