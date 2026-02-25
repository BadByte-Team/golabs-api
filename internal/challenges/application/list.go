package application

import (
	"github.com/google/uuid"

	"golabs-api/internal/challenges/domain"
)

// ListChallengesUseCase lists visible challenges for an event with optional filters.
type ListChallengesUseCase struct {
	repo domain.Repository
}

func NewListChallengesUseCase(repo domain.Repository) *ListChallengesUseCase {
	return &ListChallengesUseCase{repo: repo}
}

type ListChallengesResult struct {
	Challenge  *domain.Challenge
	SolveCount int
	FirstBlood *domain.Solve
}

func (uc *ListChallengesUseCase) Execute(
	eventID uuid.UUID,
	isAdmin bool,
	category, difficulty string,
) ([]*ListChallengesResult, error) {
	challenges, err := uc.repo.ListChallengesByEvent(eventID, !isAdmin, category, difficulty)
	if err != nil {
		return nil, err
	}

	results := make([]*ListChallengesResult, 0, len(challenges))
	for _, c := range challenges {
		count, _ := uc.repo.GetSolveCount(c.ID)
		fb, _ := uc.repo.GetFirstBlood(c.ID)
		results = append(results, &ListChallengesResult{
			Challenge:  c,
			SolveCount: count,
			FirstBlood: fb,
		})
	}
	return results, nil
}
