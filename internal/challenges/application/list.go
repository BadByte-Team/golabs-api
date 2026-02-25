// Package application contiene los casos de uso del modulo de challenges (retos CTF).
package application

import (
	"github.com/google/uuid"

	"golabs-api/internal/challenges/domain"
)

// ListChallengesUseCase lista los challenges de un evento con filtros opcionales.
// Los usuarios no admin solo ven challenges visibles (visible=true).
// Los admins ven todos los challenges incluyendo los ocultos.
type ListChallengesUseCase struct {
	repo domain.Repository
}

// NewListChallengesUseCase crea un ListChallengesUseCase con el repositorio indicado.
func NewListChallengesUseCase(repo domain.Repository) *ListChallengesUseCase {
	return &ListChallengesUseCase{repo: repo}
}

// ListChallengesResult agrupa un challenge con sus estadisticas de resoluciones.
type ListChallengesResult struct {
	Challenge  *domain.Challenge
	SolveCount int           // numero de equipos que han resuelto este challenge
	FirstBlood *domain.Solve // primer equipo en resolver el challenge (nil si nadie lo ha resuelto)
}

// Execute retorna los challenges del evento con sus estadisticas.
//
// Las estadisticas (SolveCount y FirstBlood) se obtienen en best-effort:
// si alguna falla, se omite el error para no bloquear la respuesta de toda la lista.
//
// Entrada:  eventID, isAdmin (muestra ocultos si true), category y difficulty (filtros opcionales, vacio = sin filtro).
// Salida:   slice de ListChallengesResult o error de BD.
func (uc *ListChallengesUseCase) Execute(
	eventID uuid.UUID,
	isAdmin bool,
	category, difficulty string,
) ([]*ListChallengesResult, error) {
	challenges, err := uc.repo.ListChallengesByEvent(eventID, !isAdmin, category, difficulty)
	if err != nil {
		return nil, err
	}

	results := make([]*ListChallengesResult, 0, len(challenges))
	for _, c := range challenges {
		// GetSolveCount y GetFirstBlood en best-effort: errores se ignoran intencionalmente.
		count, _ := uc.repo.GetSolveCount(c.ID)
		fb, _ := uc.repo.GetFirstBlood(c.ID)
		results = append(results, &ListChallengesResult{
			Challenge:  c,
			SolveCount: count,
			FirstBlood: fb,
		})
	}
	return results, nil
}
