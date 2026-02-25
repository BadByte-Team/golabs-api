package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

type FinishEventUseCase struct {
	repo eventdomain.Repository
}

func NewFinishEventUseCase(repo eventdomain.Repository) *FinishEventUseCase {
	return &FinishEventUseCase{repo: repo}
}

func (uc *FinishEventUseCase) Execute(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id inválido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	if err := event.Finish(); err != nil {
		return err
	}

	return uc.repo.Update(event)
}
