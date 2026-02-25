// Package application contiene los casos de uso del modulo de eventos.
package application

import (
	"time"

	eventdomain "golabs-api/internal/event/domain"
)

// CreateEventUseCase crea un nuevo evento de CTF en estado "draft".
// La apertura y el inicio del evento se realizan con use cases separados (Open, Start).
type CreateEventUseCase struct {
	repo eventdomain.Repository
}

// NewCreateEventUseCase crea un CreateEventUseCase con el repositorio indicado.
func NewCreateEventUseCase(repo eventdomain.Repository) *CreateEventUseCase {
	return &CreateEventUseCase{repo: repo}
}

// Execute valida los datos del evento, construye el agregado y lo persiste.
//
// La validacion de campos (nombre no vacio, fechas coherentes, etc.) se delega
// al constructor de dominio eventdomain.NewEvent para mantener las reglas de negocio
// en la capa de dominio.
//
// Entrada:  name, description, maxTeamSize, startsAt y endsAt del evento.
// Salida:   puntero al Event creado (en estado "draft") o error de validacion/BD.
func (uc *CreateEventUseCase) Execute(
	name string,
	description string,
	maxTeamSize int,
	startsAt time.Time,
	endsAt time.Time,
) (*eventdomain.Event, error) {

	event, err := eventdomain.NewEvent(
		name,
		description,
		maxTeamSize,
		startsAt,
		endsAt,
	)
	if err != nil {
		return nil, err
	}

	if err := uc.repo.Save(event); err != nil {
		return nil, err
	}

	return event, nil
}
