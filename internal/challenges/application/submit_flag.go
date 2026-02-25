package application

import (
	"errors"
	"fmt"

	"github.com/google/uuid"

	"golabs-api/internal/apperrors"
	challengedomain "golabs-api/internal/challenges/domain"
	eventdomain "golabs-api/internal/event/domain"
	teamdomain "golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

type SubmitFlagResult struct {
	Correct bool
	Points  int
}

type SubmitFlagUseCase struct {
	challengeRepo challengedomain.Repository
	eventRepo     eventdomain.Repository
	teamRepo      teamdomain.Repository
}

func NewSubmitFlagUseCase(
	challengeRepo challengedomain.Repository,
	eventRepo eventdomain.Repository,
	teamRepo teamdomain.Repository,
) *SubmitFlagUseCase {
	return &SubmitFlagUseCase{
		challengeRepo: challengeRepo,
		eventRepo:     eventRepo,
		teamRepo:      teamRepo,
	}
}

// Execute validates a flag submission.
//
// Security: whether the flag hash doesn't match OR the challenge has no flag set,
// the response is identical ("flag incorrecta") to avoid leaking state.
func (uc *SubmitFlagUseCase) Execute(
	challengeID uuid.UUID,
	eventID uuid.UUID,
	userID uuid.UUID,
	attempt string,
) (*SubmitFlagResult, error) {
	// 1. Verify event is running
	event, err := uc.eventRepo.GetByID(eventID)
	if err != nil {
		return nil, fmt.Errorf("%w: evento no encontrado", apperrors.ErrNotFound)
	}
	if event.Status != eventdomain.EventRunning {
		return nil, errors.New("el evento no está en curso")
	}

	// 2. Verify challenge exists and is visible
	challenge, err := uc.challengeRepo.GetChallengeByID(challengeID)
	if err != nil {
		return nil, fmt.Errorf("%w: challenge no encontrado", apperrors.ErrNotFound)
	}
	if !challenge.Visible {
		return nil, fmt.Errorf("%w: challenge no encontrado", apperrors.ErrNotFound)
	}

	// 3. Find the user's team in this event
	team, err := uc.teamRepo.GetTeamByUserAndEvent(eventID, userID)
	if err != nil {
		return nil, errors.New("debes pertenecer a un equipo para enviar una flag")
	}

	// 4. Check if this team already solved the challenge
	alreadySolved, err := uc.challengeRepo.HasTeamSolved(challengeID, team.ID)
	if err != nil {
		return nil, err
	}
	if alreadySolved {
		return nil, fmt.Errorf("%w: tu equipo ya resolvió este challenge", apperrors.ErrConflict)
	}

	// 5. Get the stored flag and compare (constant-time-ish via hash comparison)
	flag, err := uc.challengeRepo.GetFlagByChallengeID(challengeID)
	if err != nil || security.Hash(attempt) != flag.Hash {
		// Deliberately vague: wrong flag AND "no flag set" look identical
		return &SubmitFlagResult{Correct: false, Points: 0}, nil
	}

	// 6. Register the solve
	solve := challengedomain.NewSolve(challengeID, team.ID, userID)
	if err := uc.challengeRepo.SaveSolve(solve); err != nil {
		return nil, err
	}

	// 7. Add points to the team
	team.Score += challenge.Points
	team.UpdatedAt = solve.SolvedAt
	if err := uc.teamRepo.UpdateTeam(team); err != nil {
		// Solve already saved; log error but don't fail the response
		_ = err
	}

	return &SubmitFlagResult{Correct: true, Points: challenge.Points}, nil
}
