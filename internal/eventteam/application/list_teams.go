// Package application contiene los casos de uso del modulo de equipos por evento.
package application

import (
	"github.com/google/uuid"

	teamdomain "golabs-api/internal/eventteam/domain"
)

// ListTeamsByEventUseCase lista los equipos de un evento y los miembros de un equipo especifico.
// Agrupa dos operaciones relacionadas para minimizar el numero de use cases necesarios.
type ListTeamsByEventUseCase struct {
	repo teamdomain.Repository
}

// NewListTeamsByEventUseCase crea un ListTeamsByEventUseCase con el repositorio indicado.
func NewListTeamsByEventUseCase(repo teamdomain.Repository) *ListTeamsByEventUseCase {
	return &ListTeamsByEventUseCase{repo: repo}
}

// Execute retorna todos los equipos inscritos en el evento especificado.
//
// Entrada:  eventID, UUID del evento.
// Salida:   slice de EventTeam o error de BD.
func (uc *ListTeamsByEventUseCase) Execute(eventID uuid.UUID) ([]*teamdomain.EventTeam, error) {
	return uc.repo.ListTeamsByEvent(eventID)
}

// ExecuteMembers retorna todos los miembros de un equipo con sus usernames resueltos.
// Este metodo realiza un JOIN entre event_team_members y users para obtener usernames.
//
// Entrada:  teamID, UUID del equipo.
// Salida:   slice de MemberWithUsername o error de BD.
func (uc *ListTeamsByEventUseCase) ExecuteMembers(teamID uuid.UUID) ([]*teamdomain.MemberWithUsername, error) {
	return uc.repo.ListMembersWithUsername(teamID)
}
