package usecases

import (
	"errors"

	"golabs-api/internal/domain/repositories"

	"github.com/google/uuid"
)

type UpdateUserRoleUseCase struct {
	repo repositories.UserRepository
}

func NewUpdateUserRoleUseCase(repo repositories.UserRepository) *UpdateUserRoleUseCase {
	return &UpdateUserRoleUseCase{repo: repo}
}

func (uc *UpdateUserRoleUseCase) Execute(userID, role string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id inválido")
	}

	if role != "admin" && role != "user" {
		return errors.New("rol inválido")
	}

	return uc.repo.UpdateRole(userID, role)
}
