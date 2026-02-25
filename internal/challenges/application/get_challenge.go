package application

import (
	"fmt"

	"github.com/google/uuid"

	"golabs-api/internal/apperrors"
	challengedomain "golabs-api/internal/challenges/domain"
)

type GetChallengeUseCase struct {
	repo challengedomain.Repository
}

func NewGetChallengeUseCase(repo challengedomain.Repository) *GetChallengeUseCase {
	return &GetChallengeUseCase{repo: repo}
}

// Execute returns a challenge by ID.
// For non-admins, hidden challenges are treated as not found to avoid leaking existence.
func (uc *GetChallengeUseCase) Execute(id uuid.UUID, isAdmin bool) (*challengedomain.Challenge, error) {
	challenge, err := uc.repo.GetChallengeByID(id)
	if err != nil {
		return nil, fmt.Errorf("%w", apperrors.ErrNotFound)
	}

	if !isAdmin && !challenge.Visible {
		return nil, fmt.Errorf("%w: challenge not found", apperrors.ErrNotFound)
	}

	return challenge, nil
}
