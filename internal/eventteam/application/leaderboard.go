// Package application contiene los casos de uso del modulo de equipos por evento.
package application

import (
	"github.com/google/uuid"

	teamdomain "golabs-api/internal/eventteam/domain"
)

// GetLeaderboardUseCase retorna el ranking de equipos de un evento ordenado por puntaje.
// La consulta delega el ordenamiento a la BD a traves del repositorio.
type GetLeaderboardUseCase struct {
	repo teamdomain.Repository
}

// NewGetLeaderboardUseCase crea un GetLeaderboardUseCase con el repositorio indicado.
func NewGetLeaderboardUseCase(repo teamdomain.Repository) *GetLeaderboardUseCase {
	return &GetLeaderboardUseCase{repo: repo}
}

// Execute construye el ranking completo para el evento indicado.
// El rank es secuencial (1, 2, 3...) basado en el orden retornado por ListTeamsByEvent,
// que ordena por score descendente en la consulta SQL.
//
// El conteo de miembros se obtiene en un best-effort: si falla para algun equipo
// se omite el error para no bloquear la respuesta del leaderboard completo.
//
// Entrada:  eventID, UUID del evento.
// Salida:   slice de LeaderboardEntry ordenado por ranking, o error de BD.
func (uc *GetLeaderboardUseCase) Execute(eventID uuid.UUID) ([]*teamdomain.LeaderboardEntry, error) {
	teams, err := uc.repo.ListTeamsByEvent(eventID)
	if err != nil {
		return nil, err
	}

	entries := make([]*teamdomain.LeaderboardEntry, 0, len(teams))
	for i, t := range teams {
		// CountMembers en best-effort: si falla, el count queda en 0.
		count, _ := uc.repo.CountMembers(t.ID)
		entries = append(entries, &teamdomain.LeaderboardEntry{
			Rank:        i + 1,
			TeamID:      t.ID.String(),
			TeamName:    t.Name,
			Score:       t.Score,
			MemberCount: count,
		})
	}
	return entries, nil
}
