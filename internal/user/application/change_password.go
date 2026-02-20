package userapp

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	userdomain "golabs-api/internal/user/domain"
)

type ChangePasswordUseCase struct {
	repo userdomain.UserRepository
}

func NewChangePasswordUseCase(repo userdomain.UserRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo}
}

func (uc *ChangePasswordUseCase) Execute(userID, currentPassword, newPassword string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id inválido")
	}

	if currentPassword == "" || newPassword == "" {
		return errors.New("contraseñas requeridas")
	}

	if len(newPassword) < 8 {
		return errors.New("la nueva contraseña es muy corta")
	}

	user, err := uc.repo.GetByID(userID)
	if err != nil {
		return err
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(currentPassword),
	); err != nil {
		return errors.New("contraseña actual incorrecta")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	return uc.repo.UpdatePassword(userID, string(newHash))
}
