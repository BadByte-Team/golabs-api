// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"golabs-api/internal/apperrors"
	"golabs-api/internal/infrastructure/security"
	userdomain "golabs-api/internal/user/domain"
)

// CreateUserUseCase registra una nueva cuenta de usuario en el sistema.
type CreateUserUseCase struct {
	repo userdomain.UserRepository
}

// NewCreateUserUseCase crea un CreateUserUseCase con el repositorio indicado.
func NewCreateUserUseCase(repo userdomain.UserRepository) *CreateUserUseCase {
	return &CreateUserUseCase{repo: repo}
}

// Execute valida los datos, verifica unicidad y persiste el nuevo usuario.
//
// Reglas:
//   - username, email y password son obligatorios.
//   - email y username deben ser unicos en el sistema.
//   - la contrasena se hashea con bcrypt antes de persistir; nunca se almacena en texto plano.
//   - el rol inicial es siempre "user"; los admins se promueven manualmente.
//
// Entrada:  username, email y password en texto plano.
// Salida:   puntero al User creado (sin PasswordHash expuesto al handler) o error.
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

	// Verificar unicidad de email antes de continuar.
	if _, err := uc.repo.GetByEmail(email); err == nil {
		return nil, fmt.Errorf("%w: email already in use", apperrors.ErrConflict)
	}

	// Verificar unicidad de username antes de continuar.
	if _, err := uc.repo.GetByUsername(username); err == nil {
		return nil, fmt.Errorf("%w: username already in use", apperrors.ErrConflict)
	}

	// Hashear contrasena con el costo bcrypt configurado en security.BcryptCost.
	hash, err := bcrypt.GenerateFromPassword([]byte(password), security.BcryptCost)
	if err != nil {
		return nil, err
	}

	user := &userdomain.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		Role:         userdomain.RoleUser, // rol inicial siempre "user"
		Points:       0,
	}

	if err := uc.repo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}
