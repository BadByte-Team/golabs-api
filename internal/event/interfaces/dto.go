// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de eventos.
package interfaces

import "time"

// CreateEventRequest es el body de la peticion para crear un nuevo evento.
type CreateEventRequest struct {
	Name        string    `json:"name"          validate:"required,max=100"`
	Description string    `json:"description"   validate:"max=1000"`
	MaxTeamSize int       `json:"max_team_size" validate:"required,gt=0"`
	StartsAt    time.Time `json:"starts_at"     validate:"required"`
	EndsAt      time.Time `json:"ends_at"       validate:"required"`
}

// UpdateEventRequest es el body de la peticion para actualizar un evento existente.
type UpdateEventRequest struct {
	Name        string    `json:"name"          validate:"required,max=100"`
	Description string    `json:"description"   validate:"max=1000"`
	MaxTeamSize int       `json:"max_team_size" validate:"required,gt=0"`
	StartsAt    time.Time `json:"starts_at"     validate:"required"`
	EndsAt      time.Time `json:"ends_at"       validate:"required"`
}

// EventResponse es la representacion JSON de un evento para la API.
// El campo Status es el string del EventStatus (draft, open, running, finished).
type EventResponse struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxTeamSize int       `json:"max_team_size"`
	Status      string    `json:"status"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
