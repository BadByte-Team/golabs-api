package usecases

import (
	"errors"

	"golabs-api/internal/domain/repositories"

	"github.com/google/uuid"
)

type UnbanUserUseCase struct {
	repo repositories.UserRepository
}

func NewUnbanUserUseCase(repo repositories.UserRepository) *UnbanUserUseCase {
	return &UnbanUserUseCase{repo: repo}
}

func (uc *UnbanUserUseCase) Execute(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id inválido")
	}

	return uc.repo.Unban(userID)
}
