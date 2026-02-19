package usecases

import (
	"errors"

	"golabs-api/internal/domain/repositories"

	"github.com/google/uuid"
)

type UpdateUserPointsUseCase struct {
	repo repositories.UserRepository
}

func NewUpdateUserPointsUseCase(repo repositories.UserRepository) *UpdateUserPointsUseCase {
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
