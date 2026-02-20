package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

type UpdateUserRoleUseCase struct {
	repo userdomain.UserRepository
}

func NewUpdateUserRoleUseCase(repo userdomain.UserRepository) *UpdateUserRoleUseCase {
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
