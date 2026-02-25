// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	userdomain "golabs-api/internal/user/domain"
)

// GetUserByUsernameUseCase obtiene un usuario por su username exacto.
// Usado para buscar perfiles publicos donde se conoce el username completo.
type GetUserByUsernameUseCase struct {
	repo userdomain.UserRepository
}

// NewGetUserByUsernameUseCase crea un GetUserByUsernameUseCase con el repositorio indicado.
func NewGetUserByUsernameUseCase(repo userdomain.UserRepository) *GetUserByUsernameUseCase {
	return &GetUserByUsernameUseCase{repo: repo}
}

// Execute busca y retorna el usuario con el username exacto especificado.
//
// Entrada:  username, nombre de usuario exacto (case-sensitive segun configuracion de la BD).
// Salida:   puntero al User encontrado o error (apperrors.ErrNotFound si no existe).
func (uc *GetUserByUsernameUseCase) Execute(username string) (*userdomain.User, error) {
	return uc.repo.GetByUsername(username)
}
