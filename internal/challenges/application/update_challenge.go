package application

import (
	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
)

type UpdateChallengeUseCase struct {
	repo challengedomain.Repository
}

func NewUpdateChallengeUseCase(repo challengedomain.Repository) *UpdateChallengeUseCase {
	return &UpdateChallengeUseCase{repo: repo}
}

func (uc *UpdateChallengeUseCase) Execute(
	id uuid.UUID,
	title, description string,
	category challengedomain.ChallengeCategory,
	points int,
	difficulty challengedomain.ChallengeDifficulty,
) (*challengedomain.Challenge, error) {
	challenge, err := uc.repo.GetChallengeByID(id)
	if err != nil {
		return nil, err
	}

	if err := challenge.Update(title, description, category, points, difficulty); err != nil {
		return nil, err
	}

	if err := uc.repo.UpdateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
