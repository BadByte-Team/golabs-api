package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

type StartEventUseCase struct {
	repo eventdomain.Repository
}

func NewStartEventUseCase(repo eventdomain.Repository) *StartEventUseCase {
	return &StartEventUseCase{repo: repo}
}

func (uc *StartEventUseCase) Execute(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id inválido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	if err := event.Start(); err != nil {
		return err
	}

	return uc.repo.Update(event)
}
