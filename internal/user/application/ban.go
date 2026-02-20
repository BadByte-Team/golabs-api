package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

type BanUserUseCase struct {
	repo userdomain.UserRepository
}

func NewBanUserUseCase(repo userdomain.UserRepository) *BanUserUseCase {
	return &BanUserUseCase{repo: repo}
}

func (uc *BanUserUseCase) Execute(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id inválido")
	}
	return uc.repo.Ban(userID)
}
