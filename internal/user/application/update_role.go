// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

// UpdateUserRoleUseCase cambia el rol de un usuario entre "user" y "admin".
// Solo administradores pueden cambiar roles; el endpoint esta protegido por RequireRole("admin").
type UpdateUserRoleUseCase struct {
	repo userdomain.UserRepository
}

// NewUpdateUserRoleUseCase crea un UpdateUserRoleUseCase con el repositorio indicado.
func NewUpdateUserRoleUseCase(repo userdomain.UserRepository) *UpdateUserRoleUseCase {
	return &UpdateUserRoleUseCase{repo: repo}
}

// Execute asigna el nuevo rol al usuario identificado por userID.
// Los roles validos son "admin" y "user"; cualquier otro valor retorna error.
//
// Entrada:  userID (UUID en string), role ("admin" | "user").
// Salida:   error si el UUID o el rol son invalidos, o si falla la BD.
func (uc *UpdateUserRoleUseCase) Execute(userID, role string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id invalido")
	}

	if role != userdomain.RoleAdmin && role != userdomain.RoleUser {
		return errors.New("rol invalido; debe ser 'admin' o 'user'")
	}

	return uc.repo.UpdateRole(userID, role)
}
