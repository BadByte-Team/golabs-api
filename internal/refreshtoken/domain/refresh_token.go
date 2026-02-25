package refreshtokendomain

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken represents a persisted refresh token.
// The raw token value is never stored — only its SHA-256 hash.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string // SHA-256 hex
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time
}

// IsValid returns true if the token has not been revoked and has not expired.
func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}
