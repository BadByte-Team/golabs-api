package interfaces

import "time"

type CreateEventRequest struct {
	Name        string    `json:"name"          validate:"required,max=100"`
	Description string    `json:"description"   validate:"max=1000"`
	MaxTeamSize int       `json:"max_team_size" validate:"required,gt=0"`
	StartsAt    time.Time `json:"starts_at"     validate:"required"`
	EndsAt      time.Time `json:"ends_at"       validate:"required"`
}

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
