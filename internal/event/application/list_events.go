// Package application contiene los casos de uso del modulo de eventos.
package application

import eventdomain "golabs-api/internal/event/domain"

// ListEventsUseCase retorna todos los eventos del sistema sin paginar.
// La paginacion se aplica manualmente en la capa de interfaces para simplicidad,
// ya que se espera un numero reducido de eventos activos simultaneamente.
type ListEventsUseCase struct {
	repo eventdomain.Repository
}

// NewListEventsUseCase crea un ListEventsUseCase con el repositorio indicado.
func NewListEventsUseCase(repo eventdomain.Repository) *ListEventsUseCase {
	return &ListEventsUseCase{repo: repo}
}

// Execute retorna todos los eventos existentes.
//
// Salida:   slice de todos los Event o error de BD.
// TODO: implementar filtrado por estado (draft, open, running, finished) cuando sea necesario.
func (uc *ListEventsUseCase) Execute() ([]*eventdomain.Event, error) {
	return uc.repo.List()
}
