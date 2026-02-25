package application

import (
	"time"

	eventdomain "golabs-api/internal/event/domain"
)

type CreateEventUseCase struct {
	repo eventdomain.Repository
}

func NewCreateEventUseCase(repo eventdomain.Repository) *CreateEventUseCase {
	return &CreateEventUseCase{repo: repo}
}

func (uc *CreateEventUseCase) Execute(
	name string,
	description string,
	maxTeamSize int,
	startsAt time.Time,
	endsAt time.Time,
) (*eventdomain.Event, error) {

	event, err := eventdomain.NewEvent(
		name,
		description,
		maxTeamSize,
		startsAt,
		endsAt,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(event); err != nil {
		return nil, err
	}

	return event, nil
}
