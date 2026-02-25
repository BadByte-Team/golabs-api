// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

// UpdateUserPointsUseCase actualiza directamente los puntos de un usuario.
// Usado internamente por el sistema cuando un equipo resuelve un challenge.
// Solo administradores o el propio sistema (a traves de use cases internos) pueden invocar esto.
type UpdateUserPointsUseCase struct {
	repo userdomain.UserRepository
}

// NewUpdateUserPointsUseCase crea un UpdateUserPointsUseCase con el repositorio indicado.
func NewUpdateUserPointsUseCase(repo userdomain.UserRepository) *UpdateUserPointsUseCase {
	return &UpdateUserPointsUseCase{repo: repo}
}

// Execute establece los puntos del usuario identificado por userID al valor indicado.
// Los puntos deben ser mayor o igual a cero; no se permiten valores negativos.
//
// Entrada:  userID (UUID en string), points (numero entero >= 0).
// Salida:   error si el UUID es invalido, los puntos son negativos o falla la BD.
func (uc *UpdateUserPointsUseCase) Execute(userID string, points int) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id invalido")
	}

	if points < 0 {
		return errors.New("los puntos no pueden ser negativos")
	}

	return uc.repo.UpdatePoints(userID, points)
}
