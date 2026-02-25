// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

// BanUserUseCase suspende el acceso de un usuario al sistema.
// La suspension se aplica inmediatamente: el middleware LoadUser la detectara
// en la proxima peticion autenticada del usuario baneado.
type BanUserUseCase struct {
	repo userdomain.UserRepository
}

// NewBanUserUseCase crea un BanUserUseCase con el repositorio indicado.
func NewBanUserUseCase(repo userdomain.UserRepository) *BanUserUseCase {
	return &BanUserUseCase{repo: repo}
}

// Execute ban al usuario identificado por userID.
//
// Entrada:  userID, UUID del usuario a banear en formato string.
// Salida:   error si el UUID es invalido o si falla la operacion en BD.
func (uc *BanUserUseCase) Execute(userID string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id invalido")
	}
	return uc.repo.Ban(userID)
}
