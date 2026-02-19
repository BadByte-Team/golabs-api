package usecases

import (
	"errors"

	"golabs-api/internal/domain/entities"
	"golabs-api/internal/domain/repositories"

	"golang.org/x/crypto/bcrypt"
)

type CreateUserUseCase struct {
	userRepository repositories.UserRepository
}

func NewCreateUserUseCase(userRepository repositories.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{userRepository: userRepository}
}

func (uc *CreateUserUseCase) Execute(username, email, password string) (*entities.User, error) {
	if username == "" {
		return nil, errors.New("username is required")
	}
	if email == "" {
		return nil, errors.New("email is required")
	}
	if password == "" {
		return nil, errors.New("password is required")
	}

	if _, err := uc.userRepository.GetByEmail(email); err == nil {
		return nil, errors.New("email is already in use")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	if err != nil {
		return nil, err
	}

	user := &entities.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         "user",
		Points:       0,
	}

	if err := uc.userRepository.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
