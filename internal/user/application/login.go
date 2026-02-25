package userapp

import (
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"

	"golabs-api/internal/infrastructure/security"
	userdomain "golabs-api/internal/user/domain"
)

type LoginUseCase struct {
	repo userdomain.UserRepository
	jwt  *security.JWTService
}

func NewLoginUseCase(repo userdomain.UserRepository, jwt *security.JWTService) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwt: jwt}
}

// Execute accepts either an email address or a username in the `identifier` field.
func (uc *LoginUseCase) Execute(identifier, password string) (string, error) {
	var (
		user *userdomain.User
		err  error
	)

	// Determine lookup strategy: email addresses contain "@".
	if strings.Contains(identifier, "@") {
		user, err = uc.repo.GetByEmail(identifier)
	} else {
		user, err = uc.repo.GetByUsername(identifier)
	}

	if err != nil {
		// Generic error — do not reveal whether email/username exists.
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
