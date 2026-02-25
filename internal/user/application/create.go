package userapp

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"golabs-api/internal/apperrors"
	"golabs-api/internal/infrastructure/security"
	userdomain "golabs-api/internal/user/domain"
)

type CreateUserUseCase struct {
	repo userdomain.UserRepository
}

func NewCreateUserUseCase(repo userdomain.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

func (uc *CreateUserUseCase) Execute(username, email, password string) (*userdomain.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	// Enforce unique email
	if _, err := uc.repo.GetByEmail(email); err == nil {
		return nil, fmt.Errorf("%w: email already in use", apperrors.ErrConflict)
	}

	// Enforce unique username
	if _, err := uc.repo.GetByUsername(username); err == nil {
		return nil, fmt.Errorf("%w: username already in use", apperrors.ErrConflict)
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), security.BcryptCost)
	if err != nil {
		return nil, err
	}

	user := &userdomain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         userdomain.RoleUser,
		Points:       0,
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
