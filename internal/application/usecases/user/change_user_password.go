package usecases

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"golabs-api/internal/domain/repositories"
)

type ChangePasswordUseCase struct {
	repo repositories.UserRepository
}

func NewChangePasswordUseCase(repo repositories.UserRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo}
}

func (uc *ChangePasswordUseCase) Execute(
	userID string,
	currentPassword string,
	newPassword string,
) error {

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

	// Verificar contraseña actual
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(currentPassword),
	); err != nil {
		return errors.New("contraseña actual incorrecta")
	}

	// Generar nuevo hash
	newHash, err := bcrypt.GenerateFromPassword(
		[]byte(newPassword),
		12,
	)
	if err != nil {
		return err
	}

	return uc.repo.UpdatePassword(userID, string(newHash))
}
