package application

import (
	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
	eventdomain "golabs-api/internal/event/domain"
)

type CreateChallengeUseCase struct {
	challengeRepo challengedomain.Repository
	eventRepo     eventdomain.Repository
}

func NewCreateChallengeUseCase(
	challengeRepo challengedomain.Repository,
	eventRepo eventdomain.Repository,
) *CreateChallengeUseCase {
	return &CreateChallengeUseCase{
		challengeRepo: challengeRepo,
		eventRepo:     eventRepo,
	}
}

func (uc *CreateChallengeUseCase) Execute(
	eventID uuid.UUID,
	title, description string,
	category challengedomain.ChallengeCategory,
	points int,
	difficulty challengedomain.ChallengeDifficulty,
) (*challengedomain.Challenge, error) {
	// Verify the event exists
	if _, err := uc.eventRepo.GetByID(eventID); err != nil {
		return nil, err
	}

	challenge, err := challengedomain.NewChallenge(eventID, title, description, category, points, difficulty)
	if err != nil {
		return nil, err
	}

	if err := uc.challengeRepo.SaveChallenge(challenge); err != nil {
		return nil, err
	}

	return challenge, nil
}
