package interfaces

import "time"

type CreateEventRequest struct {
	Name        string    `json:"name"`
	Description string    `json:"description"`
	MaxTeamSize int       `json:"max_team_size"`
	StartsAt    time.Time `json:"starts_at"`
	EndsAt      time.Time `json:"ends_at"`
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
