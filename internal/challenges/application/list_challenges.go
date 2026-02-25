package application

import (
	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
)

type ListChallengesUseCase struct {
	repo challengedomain.Repository
}

func NewListChallengesUseCase(repo challengedomain.Repository) *ListChallengesUseCase {
	return &ListChallengesUseCase{repo: repo}
}

// Execute returns challenges for an event.
// isAdmin=true returns all (visible + hidden); isAdmin=false returns only visible ones.
func (uc *ListChallengesUseCase) Execute(eventID uuid.UUID, isAdmin bool) ([]*challengedomain.Challenge, error) {
	visibleOnly := !isAdmin
	return uc.repo.ListChallengesByEvent(eventID, visibleOnly)
}
