package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

type UnbanUserUseCase struct {
	repo userdomain.UserRepository
}

func NewUnbanUserUseCase(repo userdomain.UserRepository) *UnbanUserUseCase {
	return &UnbanUserUseCase{repo: repo}
}

func (uc *UnbanUserUseCase) Execute(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id inválido")
	}
	return uc.repo.Unban(userID)
}
