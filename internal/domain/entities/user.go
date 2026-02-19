package entities

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string
	Role         string
	Points       int
	CreatedAt    time.Time
	UpdatedAt    time.Time
}
