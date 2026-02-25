package application

import (
	"time"

	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
	"golabs-api/internal/infrastructure/security"
)

type SetFlagUseCase struct {
	repo challengedomain.Repository
}

func NewSetFlagUseCase(repo challengedomain.Repository) *SetFlagUseCase {
	return &SetFlagUseCase{repo: repo}
}

// Execute sets or replaces the flag for a challenge.
// The plaintext is hashed; only the hash is persisted.
func (uc *SetFlagUseCase) Execute(challengeID uuid.UUID, plaintext string) error {
	// Verify challenge exists
	if _, err := uc.repo.GetChallengeByID(challengeID); err != nil {
		return err
	}

	flag := &challengedomain.Flag{
		ID:          uuid.New(),
		ChallengeID: challengeID,
		Hash:        security.Hash(plaintext),
		CreatedAt:   time.Now(),
	}

	return uc.repo.UpsertFlag(flag)
}
