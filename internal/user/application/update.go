// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

// UpdateUserUseCase modifica los datos de perfil editables de un usuario (username, email).
// La contrasena y el rol se manejan con use cases dedicados por separado.
type UpdateUserUseCase struct {
	repo userdomain.UserRepository
}

// NewUpdateUserUseCase crea un UpdateUserUseCase con el repositorio indicado.
func NewUpdateUserUseCase(repo userdomain.UserRepository) *UpdateUserUseCase {
	return &UpdateUserUseCase{repo: repo}
}

// Execute actualiza username y/o email del usuario. Los campos vacios se ignoran
// (patch semantics: solo se actualizan los campos que vienen con valor).
//
// Entrada:  id (UUID en string), username y email nuevos (pueden ser cadenas vacias para no cambiarlos).
// Salida:   puntero al User actualizado o error.
func (uc *UpdateUserUseCase) Execute(id, username, email string) (*userdomain.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("id invalido")
	}

	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	// Solo actualizar los campos que se proporcionan (semantica de PATCH).
	if username != "" {
		user.Username = username
	}

	if email != "" {
		user.Email = email
	}

	if err := uc.repo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}
