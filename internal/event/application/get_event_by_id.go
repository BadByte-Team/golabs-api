package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

type GetEventByIDUseCase struct {
	repo eventdomain.Repository
}

func NewGetEventByIDUseCase(repo eventdomain.Repository) *GetEventByIDUseCase {
	return &GetEventByIDUseCase{repo: repo}
}

func (uc *GetEventByIDUseCase) Execute(id string) (*eventdomain.Event, error) {
	eventID, err := uuid.Parse(id)
	if err != nil {
		return nil, errors.New("id inválido")
	}

	return uc.repo.GetByID(eventID)
}
