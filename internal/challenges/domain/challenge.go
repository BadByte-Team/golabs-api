package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

/* ── Category ── */

type ChallengeCategory string

const (
	CategoryPwn       ChallengeCategory = "pwn"
	CategoryWeb       ChallengeCategory = "web"
	CategoryCrypto    ChallengeCategory = "crypto"
	CategoryForensics ChallengeCategory = "forensics"
	CategoryReverse   ChallengeCategory = "reverse"
	CategoryOSINT     ChallengeCategory = "osint"
	CategoryMisc      ChallengeCategory = "misc"
)

/* ── Difficulty ── */

type ChallengeDifficulty string

const (
	DifficultyEasy   ChallengeDifficulty = "easy"
	DifficultyMedium ChallengeDifficulty = "medium"
	DifficultyHard   ChallengeDifficulty = "hard"
)

/* ── Challenge ── */

type Challenge struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	Title       string
	Description string
	Category    ChallengeCategory
	Points      int
	Difficulty  ChallengeDifficulty
	Visible     bool
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

func NewChallenge(
	eventID uuid.UUID,
	title, description string,
	category ChallengeCategory,
	points int,
	difficulty ChallengeDifficulty,
) (*Challenge, error) {
	if eventID == uuid.Nil {
		return nil, errors.New("eventID requerido")
	}
	if title == "" {
		return nil, errors.New("title requerido")
	}
	if points < 0 {
		return nil, errors.New("points no puede ser negativo")
	}
	if category == "" {
		return nil, errors.New("category requerida")
	}
	if difficulty == "" {
		difficulty = DifficultyMedium
	}

	now := time.Now()
	return &Challenge{
		ID:          uuid.New(),
		EventID:     eventID,
		Title:       title,
		Description: description,
		Category:    category,
		Points:      points,
		Difficulty:  difficulty,
		Visible:     false, // oculto por defecto hasta que el admin lo publique
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update changes the editable fields of the challenge.
func (c *Challenge) Update(title, description string, category ChallengeCategory, points int, difficulty ChallengeDifficulty) error {
	if title == "" {
		return errors.New("title requerido")
	}
	if points < 0 {
		return errors.New("points no puede ser negativo")
	}
	c.Title = title
	c.Description = description
	c.Category = category
	c.Points = points
	c.Difficulty = difficulty
	c.UpdatedAt = time.Now()
	return nil
}

// Publish makes the challenge visible to participants.
func (c *Challenge) Publish() {
	c.Visible = true
	c.UpdatedAt = time.Now()
}

// Unpublish hides the challenge from participants.
func (c *Challenge) Unpublish() {
	c.Visible = false
	c.UpdatedAt = time.Now()
}

/* ── Flag ── */

// Flag stores the SHA-256 hash of the real flag value.
// The plain-text value is NEVER persisted.
type Flag struct {
	ID          uuid.UUID
	ChallengeID uuid.UUID
	Hash        string // hex-encoded SHA-256
	CreatedAt   time.Time
}

/* ── Solve ── */

// Solve is an immutable record created when a team submits the correct flag.
type Solve struct {
	ID          uuid.UUID
	ChallengeID uuid.UUID
	EventTeamID uuid.UUID
	UserID      uuid.UUID
	SolvedAt    time.Time
}

func NewSolve(challengeID, teamID, userID uuid.UUID) *Solve {
	return &Solve{
		ID:          uuid.New(),
		ChallengeID: challengeID,
		EventTeamID: teamID,
		UserID:      userID,
		SolvedAt:    time.Now(),
	}
}
