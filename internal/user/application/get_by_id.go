// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

// GetUserByIDUseCase obtiene un usuario por su UUID.
// Usado por handlers de perfil donde se conoce el ID directamente del token o de la URL.
type GetUserByIDUseCase struct {
	repo userdomain.UserRepository
}

// NewGetUserByIDUseCase crea un GetUserByIDUseCase con el repositorio indicado.
func NewGetUserByIDUseCase(repo userdomain.UserRepository) *GetUserByIDUseCase {
	return &GetUserByIDUseCase{repo: repo}
}

// Execute busca y retorna el usuario con el ID especificado.
//
// Entrada:  id, UUID del usuario en formato string.
// Salida:   puntero al User encontrado o error (apperrors.ErrNotFound si no existe).
func (uc *GetUserByIDUseCase) Execute(id string) (*userdomain.User, error) {
	if _, err := uuid.Parse(id); err != nil {
		return nil, errors.New("id invalido")
	}

	user, err := uc.repo.GetByID(id)
	if err != nil {
		return nil, err
	}

	return user, nil
}
