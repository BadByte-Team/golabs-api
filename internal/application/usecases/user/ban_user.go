package usecases

import (
	"errors"

	"golabs-api/internal/domain/repositories"

	"github.com/google/uuid"
)

type BanUserUseCase struct {
	repo repositories.UserRepository
}

func NewBanUserUseCase(repo repositories.UserRepository) *BanUserUseCase {
	return &BanUserUseCase{repo: repo}
}

func (uc *BanUserUseCase) Execute(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id inválido")
	}

	return uc.repo.Ban(userID)
}
