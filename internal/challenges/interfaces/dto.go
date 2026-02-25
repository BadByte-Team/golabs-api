package interfaces

import "time"

/* ── Challenge ─────────────────────────────────────────────────────────── */

type CreateChallengeRequest struct {
	Title       string `json:"title"       validate:"required,max=120"`
	Description string `json:"description" validate:"required,max=2000"`
	Category    string `json:"category"    validate:"required,oneof=web pwn rev crypto forensics misc"`
	Points      int    `json:"points"      validate:"required,gt=0"`
	Difficulty  string `json:"difficulty"  validate:"required,oneof=easy medium hard"`
}

type UpdateChallengeRequest struct {
	Title       string `json:"title"       validate:"required,max=120"`
	Description string `json:"description" validate:"required,max=2000"`
	Category    string `json:"category"    validate:"required,oneof=web pwn rev crypto forensics misc"`
	Points      int    `json:"points"      validate:"required,gt=0"`
	Difficulty  string `json:"difficulty"  validate:"required,oneof=easy medium hard"`
}

type ChallengeResponse struct {
	ID               string    `json:"id"`
	EventID          string    `json:"event_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Category         string    `json:"category"`
	Points           int       `json:"points"`
	Difficulty       string    `json:"difficulty"`
	Visible          bool      `json:"visible"`
	SolveCount       int       `json:"solve_count"`
	FirstBloodTeamID *string   `json:"first_blood_team_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

/* ── Flag ──────────────────────────────────────────────────────────────── */

// SetFlagRequest is only sent by admins; the plain-text is hashed server-side.
type SetFlagRequest struct {
	Flag string `json:"flag" validate:"required,min=1"`
}

/* ── Submit ────────────────────────────────────────────────────────────── */

type SubmitFlagRequest struct {
	Flag string `json:"flag" validate:"required"`
}

type SubmitFlagResponse struct {
	Correct bool `json:"correct"`
	Points  int  `json:"points,omitempty"`
}

/* ── Solve ─────────────────────────────────────────────────────────────── */

type SolveResponse struct {
	ChallengeID string    `json:"challenge_id"`
	TeamID      string    `json:"event_team_id"`
	UserID      string    `json:"user_id"`
	SolvedAt    time.Time `json:"solved_at"`
}
