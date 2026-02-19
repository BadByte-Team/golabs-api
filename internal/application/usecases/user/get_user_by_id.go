package usecases

import (
	"errors"

	"github.com/google/uuid"

	"golabs-api/internal/domain/entities"
	"golabs-api/internal/domain/repositories"
)

type GetUserByIDUseCase struct {
	repository repositories.UserRepository
}

func NewGetUserByIDUseCase(repository repositories.UserRepository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{repository: repository}
}

func (uc *GetUserByIDUseCase) Execute(id string) (*entities.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("id inválido")
	}

	user, err := uc.repository.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
