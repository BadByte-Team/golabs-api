// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de equipos.
package interfaces

import "time"

// CreateTeamRequest es el body para crear un nuevo equipo en un evento.
type CreateTeamRequest struct {
	Name string `json:"name" validate:"required,min=2,max=40"`
}

// JoinTeamRequest es el body para unirse a un equipo existente.
// Requiere el nombre exacto del equipo y el join secret compartido por el owner.
type JoinTeamRequest struct {
	TeamName   string `json:"team_name"   validate:"required"`
	JoinSecret string `json:"join_secret" validate:"required"`
}

// EventTeamResponse es la representacion JSON de un equipo para la API.
type EventTeamResponse struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	Name        string `json:"name"`
	Score       int    `json:"score"`
	MemberCount int    `json:"member_count,omitempty"`
}

// CreateTeamResponse extiende EventTeamResponse con el join secret en texto plano.
// El join secret solo se incluye en la respuesta de creacion; despues no se puede recuperar.
type CreateTeamResponse struct {
	EventTeamResponse
	JoinSecret string `json:"join_secret"`
}

// EventTeamMemberResponse representa a un miembro de equipo con su username resuelto.
type EventTeamMemberResponse struct {
	UserID   string    `json:"user_id"`
	Username string    `json:"username"`
	Role     string    `json:"role"`
	JoinedAt time.Time `json:"joined_at"`
}
