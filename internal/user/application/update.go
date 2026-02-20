package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

type UpdateUserUseCase struct {
	repo userdomain.UserRepository
}

func NewUpdateUserUseCase(repo userdomain.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

func (uc *UpdateUserUseCase) Execute(id, username, email string) (*userdomain.User, error) {
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
