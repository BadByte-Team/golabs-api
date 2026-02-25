// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

// UnbanUserUseCase reactiva el acceso de un usuario baneado al sistema.
// Solo administradores pueden desbanear usuarios.
type UnbanUserUseCase struct {
	repo userdomain.UserRepository
}

// NewUnbanUserUseCase crea un UnbanUserUseCase con el repositorio indicado.
func NewUnbanUserUseCase(repo userdomain.UserRepository) *UnbanUserUseCase {
	return &UnbanUserUseCase{repo: repo}
}

// Execute levanta el baneo del usuario identificado por userID.
//
// Entrada:  userID, UUID del usuario a desbanear en formato string.
// Salida:   error si el UUID es invalido o si falla la operacion en BD.
func (uc *UnbanUserUseCase) Execute(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id invalido")
	}
	return uc.repo.Unban(userID)
}
