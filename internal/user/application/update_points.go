package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

type UpdateUserPointsUseCase struct {
	repo userdomain.UserRepository
}

func NewUpdateUserPointsUseCase(repo userdomain.UserRepository) *UpdateUserPointsUseCase {
	return &UpdateUserPointsUseCase{repo: repo}
}

func (uc *UpdateUserPointsUseCase) Execute(userID string, points int) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id inválido")
	}

	if points < 0 {
		return errors.New("los puntos no pueden ser negativos")
	}

	return uc.repo.UpdatePoints(userID, points)
}
