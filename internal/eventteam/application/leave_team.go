// Package application contiene los casos de uso del modulo de equipos por evento.
package application

import (
	"errors"

	"github.com/google/uuid"

	"golabs-api/internal/eventteam/domain"
)

// LeaveTeamUseCase elimina a un usuario de su equipo.
// El owner no puede abandonar el equipo si hay otros miembros;
// primero debe transferir la propiedad o esperar a que los miembros salgan.
type LeaveTeamUseCase struct {
	teamsRepo domain.Repository
}

// NewLeaveTeamUseCase crea un LeaveTeamUseCase con el repositorio indicado.
func NewLeaveTeamUseCase(
	teamsRepo domain.Repository,
) *LeaveTeamUseCase {
	return &LeaveTeamUseCase{teamsRepo: teamsRepo}
}

// Execute elimina al usuario del equipo, con validacion de rol.
//
// Reglas:
//   - El usuario debe ser miembro del equipo.
//   - El owner solo puede salir si es el unico miembro (equipo vacio despues de salir).
//
// Entrada:  teamID, userID (UUIDs del equipo y el usuario).
// Salida:   error si el usuario no pertenece al equipo, viola las reglas de owner o falla la BD.
func (uc *LeaveTeamUseCase) Execute(
	teamID uuid.UUID,
	userID uuid.UUID,
) error {

	members, err := uc.teamsRepo.ListMembers(teamID)
	if err != nil {
		return err
	}

	// Identificar el rol del usuario en el equipo.
	var role domain.TeamRole
	for _, m := range members {
		if m.UserID == userID {
			role = m.Role
			break
		}
	}

	if role == "" {
		return errors.New("usuario no pertenece al equipo")
	}

	// El owner no puede salir si hay otros miembros activos.
	if role == domain.TeamOwner && len(members) > 1 {
		return errors.New("owner no puede salir si hay otros miembros")
	}

	return uc.teamsRepo.RemoveMember(teamID, userID)
}
