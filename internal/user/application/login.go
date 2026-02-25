package userapp

import (
	"errors"
	"strings"

	"github.com/google/uuid"
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
// Returns the signed access token and the user's UUID (needed to issue a refresh token).
func (uc *LoginUseCase) Execute(identifier, password string) (accessToken string, userID uuid.UUID, err error) {
	var user *userdomain.User

	// Determine lookup strategy: email addresses contain "@".
	if strings.Contains(identifier, "@") {
		user, err = uc.repo.GetByEmail(identifier)
	} else {
		user, err = uc.repo.GetByUsername(identifier)
	}

	if err != nil {
		// Generic error — do not reveal whether email/username exists.
		return "", uuid.Nil, errors.New("credenciales inválidas")
	}

	if user.Banned {
		return "", uuid.Nil, errors.New("usuario baneado")
	}

	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", uuid.Nil, errors.New("credenciales inválidas")
	}

	token, err := uc.jwt.Generate(user.ID.String(), user.Role)
	if err != nil {
		return "", uuid.Nil, err
	}

	return token, user.ID, nil
}
