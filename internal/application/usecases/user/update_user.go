package usecases

import (
	"errors"

	"github.com/google/uuid"

	"golabs-api/internal/domain/entities"
	"golabs-api/internal/domain/repositories"
)

type UpdateUserUseCase struct {
	repo repositories.UserRepository
}

func NewUpdateUserUseCase(repo repositories.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Execute(
	id string,
	username string,
	email string,
) (*entities.User, error) {

	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("id inválido")
	}

	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	if username != "" {
		user.Username = username
	}

	if email != "" {
		user.Email = email
	}

	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
