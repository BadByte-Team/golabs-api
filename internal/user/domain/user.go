package userdomain

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

	Banned   bool
	BannedAt *time.Time

	CreatedAt time.Time
	UpdatedAt time.Time
}
