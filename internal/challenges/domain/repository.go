package domain

import "github.com/google/uuid"

type Repository interface {
	// ── Challenges ──────────────────────────────────────────────────────────

	SaveChallenge(c *Challenge) error
	UpdateChallenge(c *Challenge) error
	GetChallengeByID(id uuid.UUID) (*Challenge, error)
	// visibleOnly=true for participants; false for admins
	// category and difficulty are optional filters (empty = no filter)
	ListChallengesByEvent(eventID uuid.UUID, visibleOnly bool, category, difficulty string) ([]*Challenge, error)

	// ── Flags ───────────────────────────────────────────────────────────────

	// UpsertFlag inserts or replaces the flag for a challenge (1 flag per challenge).
	UpsertFlag(f *Flag) error
	GetFlagByChallengeID(challengeID uuid.UUID) (*Flag, error)

	// ── Solves ──────────────────────────────────────────────────────────────

	SaveSolve(s *Solve) error
	HasTeamSolved(challengeID, teamID uuid.UUID) (bool, error)
	ListSolvesByChallenge(challengeID uuid.UUID) ([]*Solve, error)
	ListSolvesByTeam(teamID uuid.UUID) ([]*Solve, error)

	// GetSolveCount returns how many teams have solved this challenge.
	GetSolveCount(challengeID uuid.UUID) (int, error)
	// GetFirstBlood returns the first solve for a challenge (nil if unsolved).
	GetFirstBlood(challengeID uuid.UUID) (*Solve, error)
}
