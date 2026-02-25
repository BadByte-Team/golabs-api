package application

import (
	"github.com/google/uuid"

	teamdomain "golabs-api/internal/eventteam/domain"
)

// GetLeaderboardUseCase returns teams in an event ranked by score.
type GetLeaderboardUseCase struct {
	repo teamdomain.Repository
}

func NewGetLeaderboardUseCase(repo teamdomain.Repository) *GetLeaderboardUseCase {
	return &GetLeaderboardUseCase{repo: repo}
}

func (uc *GetLeaderboardUseCase) Execute(eventID uuid.UUID) ([]*teamdomain.LeaderboardEntry, error) {
	teams, err := uc.repo.ListTeamsByEvent(eventID)
	if err != nil {
		return nil, err
	}

	entries := make([]*teamdomain.LeaderboardEntry, 0, len(teams))
	for i, t := range teams {
		count, _ := uc.repo.CountMembers(t.ID) // best-effort; ignore error
		entries = append(entries, &teamdomain.LeaderboardEntry{
			Rank:        i + 1,
			TeamID:      t.ID.String(),
			TeamName:    t.Name,
			Score:       t.Score,
			MemberCount: count,
		})
	}
	return entries, nil
}
