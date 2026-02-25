package interfaces

import "time"

/* ── Challenge ─────────────────────────────────────────────────────────── */

type CreateChallengeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Points      int    `json:"points"`
	Difficulty  string `json:"difficulty"`
}

type UpdateChallengeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Points      int    `json:"points"`
	Difficulty  string `json:"difficulty"`
}

type ChallengeResponse struct {
	ID          string    `json:"id"`
	EventID     string    `json:"event_id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	Category    string    `json:"category"`
	Points      int       `json:"points"`
	Difficulty  string    `json:"difficulty"`
	Visible     bool      `json:"visible"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

/* ── Flag ──────────────────────────────────────────────────────────────── */

// SetFlagRequest is only sent by admins; the plain-text is hashed server-side.
type SetFlagRequest struct {
	Flag string `json:"flag"`
}

/* ── Submit ────────────────────────────────────────────────────────────── */

type SubmitFlagRequest struct {
	Flag string `json:"flag"`
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
