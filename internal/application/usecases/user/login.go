package usecases

import (
	"errors"

	"golabs-api/internal/domain/repositories"
	"golabs-api/internal/infrastructure/security"

	"golang.org/x/crypto/bcrypt"
)

type LoginUseCase struct {
	repo repositories.UserRepository
	jwt  *security.JWTService
}

func NewLoginUseCase(
	repo repositories.UserRepository,
	jwt *security.JWTService,
) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwt: jwt}
}

func (uc *LoginUseCase) Execute(email, password string) (string, error) {
	user, err := uc.repo.GetByEmail(email)
	if err != nil {
		return "", errors.New("credenciales inválidas")
	}

	if user.Banned {
		return "", errors.New("usuario baneado")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", errors.New("credenciales inválidas")
	}

	return uc.jwt.Generate(user.ID.String(), user.Role)
}
