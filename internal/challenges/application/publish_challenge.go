package application

import (
	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
)

type PublishChallengeUseCase struct {
	repo challengedomain.Repository
}

func NewPublishChallengeUseCase(repo challengedomain.Repository) *PublishChallengeUseCase {
	return &PublishChallengeUseCase{repo: repo}
}

// Execute publishes (visible=true) or unpublishes (visible=false) a challenge.
func (uc *PublishChallengeUseCase) Execute(id uuid.UUID, publish bool) (*challengedomain.Challenge, error) {
	challenge, err := uc.repo.GetChallengeByID(id)
	if err != nil {
		return nil, err
	}

	if publish {
		challenge.Publish()
	} else {
		challenge.Unpublish()
	}

	if err := uc.repo.UpdateChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
