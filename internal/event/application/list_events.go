package application

import eventdomain "golabs-api/internal/event/domain"

type ListEventsUseCase struct {
	repo eventdomain.Repository
}

func NewListEventsUseCase(repo eventdomain.Repository) *ListEventsUseCase {
	return &ListEventsUseCase{repo: repo}
}

func (uc *ListEventsUseCase) Execute() ([]*eventdomain.Event, error) {
	return uc.repo.List()
}
