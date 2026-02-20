package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

type GetUserByIDUseCase struct {
	repo userdomain.UserRepository
}

func NewGetUserByIDUseCase(repo userdomain.UserRepository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{repo: repo}
}

func (uc *GetUserByIDUseCase) Execute(id string) (*userdomain.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("id inválido")
	}

	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
