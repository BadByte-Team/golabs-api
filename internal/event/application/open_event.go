package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

type OpenEventUseCase struct {
	repo eventdomain.Repository
}

func NewOpenEventUseCase(repo eventdomain.Repository) *OpenEventUseCase {
	return &OpenEventUseCase{repo: repo}
}

func (uc *OpenEventUseCase) Execute(id string) error {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return errors.New("id inválido")
	}

	event, err := uc.repo.GetByID(eventID)
	if err != nil {
		return err
	}

	if err := event.Open(); err != nil {
		return err
	}

	return uc.repo.Update(event)
}
