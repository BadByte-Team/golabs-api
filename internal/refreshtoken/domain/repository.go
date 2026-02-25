package refreshtokendomain

import (
	"context"

	"github.com/google/uuid"
)

// RefreshTokenRepository defines persistence operations for refresh tokens.
type RefreshTokenRepository interface {
	// Save persists a new refresh token.
	Save(ctx context.Context, rt *RefreshToken) error

	// GetByTokenHash returns the refresh token matching the given SHA-256 hash.
	GetByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)

	// Revoke marks a single token as revoked.
	Revoke(ctx context.Context, id uuid.UUID) error

	// RevokeAllForUser revokes every active token for a user (e.g., password change).
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
