// Package interfaces implementa los handlers HTTP y el registro de rutas del modulo de challenges.
package interfaces

import "time"

// CreateChallengeRequest es el body para crear un nuevo challenge en un evento (admin only).
type CreateChallengeRequest struct {
	Title       string `json:"title"       validate:"required,max=120"`
	Description string `json:"description" validate:"required,max=2000"`
	Category    string `json:"category"    validate:"required,oneof=web pwn rev crypto forensics misc"`
	Points      int    `json:"points"      validate:"required,gt=0"`
	Difficulty  string `json:"difficulty"  validate:"required,oneof=easy medium hard"`
	FileURL     string `json:"file_url"    validate:"omitempty,url,max=512"`
}

// UpdateChallengeRequest es el body para actualizar un challenge existente (admin only).
type UpdateChallengeRequest struct {
	Title       string `json:"title"       validate:"required,max=120"`
	Description string `json:"description" validate:"required,max=2000"`
	Category    string `json:"category"    validate:"required,oneof=web pwn rev crypto forensics misc"`
	Points      int    `json:"points"      validate:"required,gt=0"`
	Difficulty  string `json:"difficulty"  validate:"required,oneof=easy medium hard"`
	FileURL     string `json:"file_url"    validate:"omitempty,url,max=512"`
}

// ChallengeResponse es la representacion JSON de un challenge para la API.
// SolveCount y FirstBloodTeamID son nulos cuando se retorna un challenge sin estadisticas.
type ChallengeResponse struct {
	ID               string    `json:"id"`
	EventID          string    `json:"event_id"`
	Title            string    `json:"title"`
	Description      string    `json:"description"`
	Category         string    `json:"category"`
	Points           int       `json:"points"`
	Difficulty       string    `json:"difficulty"`
	FileURL          string    `json:"file_url,omitempty"`
	Visible          bool      `json:"visible"`
	SolveCount       int       `json:"solve_count"`
	FirstBloodTeamID *string   `json:"first_blood_team_id,omitempty"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// SetFlagRequest es el body para establecer la flag de un challenge (admin only).
// El texto plano se hashea en el servidor; nunca se almacena en claro.
type SetFlagRequest struct {
	Flag string `json:"flag" validate:"required,min=1"`
}

// SubmitFlagRequest es el body para enviar una flag (participants).
type SubmitFlagRequest struct {
	Flag string `json:"flag" validate:"required"`
}

// SubmitFlagResponse indica si la flag fue correcta y cuantos puntos se otorgaron.
// Points es 0 si la flag fue incorrecta o si el equipo ya habia resuelto el challenge.
type SubmitFlagResponse struct {
	Correct bool `json:"correct"`
	Points  int  `json:"points,omitempty"`
}

// SolveResponse es la representacion JSON de un solve (resolucion de challenge por un equipo).
type SolveResponse struct {
	ChallengeID string    `json:"challenge_id"`
	TeamID      string    `json:"event_team_id"`
	UserID      string    `json:"user_id"`
	SolvedAt    time.Time `json:"solved_at"`
}
