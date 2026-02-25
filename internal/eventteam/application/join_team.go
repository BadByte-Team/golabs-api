// Package application contiene los casos de uso del modulo de equipos por evento.
package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
	"golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

// JoinTeamUseCase une a un usuario a un equipo existente en un evento.
// Requiere que el usuario conozca el nombre del equipo y su join secret.
type JoinTeamUseCase struct {
	eventsRepo eventdomain.Repository
	teamsRepo  domain.Repository
}

// NewJoinTeamUseCase crea un JoinTeamUseCase con los repositorios indicados.
func NewJoinTeamUseCase(
	eventsRepo eventdomain.Repository,
	teamsRepo domain.Repository,
) *JoinTeamUseCase {
	return &JoinTeamUseCase{
		eventsRepo: eventsRepo,
		teamsRepo:  teamsRepo,
	}
}

// Execute valida y ejecuta la union de un usuario a un equipo.
//
// Reglas validadas en este orden:
//  1. El evento debe existir y estar en estado "open".
//  2. El usuario no puede pertenecer ya a otro equipo en el mismo evento.
//  3. El equipo debe existir por ese nombre en ese evento.
//  4. El join secret proporcionado debe coincidir con el hash almacenado.
//  5. El equipo no debe estar lleno (segun maxTeamSize del evento).
//
// Los errores de join secret incorrecto y equipo no encontrado se unifican en un mensaje
// generico para no revelar cuales equipos existen en el evento.
//
// Entrada:  eventID, userID (UUIDs), teamName y joinSecret en texto plano.
// Salida:   error si alguna validacion falla o si falla la BD.
func (uc *JoinTeamUseCase) Execute(
	eventID uuid.UUID,
	userID uuid.UUID,
	teamName string,
	joinSecret string,
) error {

	event, err := uc.eventsRepo.GetByID(eventID)
	if err != nil {
		return errors.New("evento no encontrado")
	}

	if !event.IsOpen() {
		return errors.New("evento no esta abierto")
	}

	inEvent, err := uc.teamsRepo.IsUserInEvent(eventID, userID)
	if err != nil {
		return err
	}
	if inEvent {
		return errors.New("usuario ya pertenece a un equipo en este evento")
	}

	team, err := uc.teamsRepo.GetTeamByName(eventID, teamName)
	if err != nil {
		// Error generico para no revelar si el equipo existe.
		return errors.New("no se pudo unir al equipo")
	}

	// Verificar join secret comparando hashes (nunca texto plano).
	if security.HashJoinSecret(joinSecret) != team.JoinSecretHash {
		return errors.New("no se pudo unir al equipo")
	}

	count, err := uc.teamsRepo.CountMembers(team.ID)
	if err != nil {
		return err
	}
	if count >= event.MaxTeamSize {
		return errors.New("equipo lleno")
	}

	member := &domain.EventTeamMember{
		EventTeamID: team.ID,
		UserID:      userID,
		Role:        domain.TeamMember,
	}

	return uc.teamsRepo.AddMember(member)
}
