package refreshtokeninfra

import (
	"context"
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
)

// MySQLRefreshTokenRepository implements RefreshTokenRepository using MySQL/MariaDB.
type MySQLRefreshTokenRepository struct {
	db *sql.DB
}

// New creates a new MySQLRefreshTokenRepository.
func New(db *sql.DB) *MySQLRefreshTokenRepository {
	return &MySQLRefreshTokenRepository{db: db}
}

// HashToken returns the SHA-256 hex of a raw token string.
func HashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

// Save inserts a new refresh token row.
func (r *MySQLRefreshTokenRepository) Save(ctx context.Context, rt *refreshtokendomain.RefreshToken) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO refresh_tokens (id, user_id, token_hash, expires_at, created_at)
		 VALUES (?, ?, ?, ?, ?)`,
		rt.ID.String(), rt.UserID.String(),
		rt.TokenHash, rt.ExpiresAt, rt.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("save refresh token: %w", err)
	}
	return nil
}

// GetByTokenHash fetches a refresh token by its SHA-256 hash.
func (r *MySQLRefreshTokenRepository) GetByTokenHash(ctx context.Context, hash string) (*refreshtokendomain.RefreshToken, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
		 FROM refresh_tokens WHERE token_hash = ?`,
		hash,
	)
	return scanRefreshToken(row)
}

// Revoke sets revoked_at = NOW() for the given token ID.
func (r *MySQLRefreshTokenRepository) Revoke(ctx context.Context, id uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = ? WHERE id = ?`,
		time.Now(), id.String(),
	)
	if err != nil {
		return fmt.Errorf("revoke refresh token: %w", err)
	}
	return nil
}

// RevokeAllForUser revokes all non-revoked tokens for a user.
func (r *MySQLRefreshTokenRepository) RevokeAllForUser(ctx context.Context, userID uuid.UUID) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE refresh_tokens SET revoked_at = ?
		 WHERE user_id = ? AND revoked_at IS NULL`,
		time.Now(), userID.String(),
	)
	if err != nil {
		return fmt.Errorf("revoke all tokens for user: %w", err)
	}
	return nil
}

// ── helpers ──────────────────────────────────────────────────────────────────

func scanRefreshToken(row *sql.Row) (*refreshtokendomain.RefreshToken, error) {
	var rt refreshtokendomain.RefreshToken
	var idStr, userIDStr string
	var revokedAt sql.NullTime

	err := row.Scan(
		&idStr, &userIDStr, &rt.TokenHash,
		&rt.ExpiresAt, &rt.CreatedAt, &revokedAt,
	)
	if err != nil {
		return nil, err
	}

	rt.ID, err = uuid.Parse(idStr)
	if err != nil {
		return nil, fmt.Errorf("parse refresh token id: %w", err)
	}
	rt.UserID, err = uuid.Parse(userIDStr)
	if err != nil {
		return nil, fmt.Errorf("parse user id: %w", err)
	}
	if revokedAt.Valid {
		rt.RevokedAt = &revokedAt.Time
	}

	return &rt, nil
}
