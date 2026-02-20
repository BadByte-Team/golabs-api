package userapp

import (
	"errors"

	"golang.org/x/crypto/bcrypt"

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

	if _, err := uc.repo.GetByEmail(email); err == nil {
		return nil, errors.New("email is already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := &userdomain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
		Points:       0,
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
