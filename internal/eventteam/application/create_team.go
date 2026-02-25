// Package application contiene los casos de uso del modulo de equipos por evento.
package application

import (
	"errors"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
	"golabs-api/internal/eventteam/domain"
	"golabs-api/internal/infrastructure/security"
)

// CreateTeamResult contiene el equipo creado y el join secret en texto plano.
// El join secret solo se retorna en esta operacion; despues solo se almacena su hash.
type CreateTeamResult struct {
	Team       *domain.EventTeam
	JoinSecret string // valor crudo, UNA SOLA VEZ
}

// CreateTeamUseCase crea un nuevo equipo en un evento y agrega al creador como owner.
// Solo se puede crear un equipo si el evento esta en estado "open".
type CreateTeamUseCase struct {
	eventsRepo eventdomain.Repository
	teamsRepo  domain.Repository
}

// NewCreateTeamUseCase crea un CreateTeamUseCase con los repositorios indicados.
func NewCreateTeamUseCase(
	eventsRepo eventdomain.Repository,
	teamsRepo domain.Repository,
) *CreateTeamUseCase {
	return &CreateTeamUseCase{
		eventsRepo: eventsRepo,
		teamsRepo:  teamsRepo,
	}
}

// Execute crea el equipo y lo persiste con el join secret hasheado.
//
// Flujo:
//  1. Verificar que el evento exista y este en estado "open".
//  2. Generar un join secret aleatorio (10 bytes, codificado en base64).
//  3. Hashear el join secret y crear el equipo con el hash.
//  4. Agregar al creador como miembro con rol "owner".
//
// Entrada:  eventID, ownerID (UUIDs del evento y el creador), teamName.
// Salida:   CreateTeamResult con el equipo y el join secret crudo, o error.
func (uc *CreateTeamUseCase) Execute(
	eventID uuid.UUID,
	ownerID uuid.UUID,
	teamName string,
) (*CreateTeamResult, error) {

	event, err := uc.eventsRepo.GetByID(eventID)
	if err != nil {
		return nil, errors.New("evento no encontrado")
	}

	if !event.IsOpen() {
		return nil, errors.New("evento no esta abierto")
	}

	// Generar join secret aleatorio e inmediantamente hashear para la BD.
	joinSecret, err := security.GenerateJoinSecret(10)
	if err != nil {
		return nil, err
	}

	hash := security.HashJoinSecret(joinSecret)

	team, err := domain.NewEventTeam(eventID, teamName, hash)
	if err != nil {
		return nil, err
	}

	if err := uc.teamsRepo.SaveTeam(team); err != nil {
		return nil, err
	}

	member := &domain.EventTeamMember{
		EventTeamID: team.ID,
		UserID:      ownerID,
		Role:        domain.TeamOwner,
	}

	if err := uc.teamsRepo.AddMember(member); err != nil {
		return nil, err
	}

	return &CreateTeamResult{
		Team:       team,
		JoinSecret: joinSecret, // retornar valor crudo al handler para mostrarlo al usuario
	}, nil
}
