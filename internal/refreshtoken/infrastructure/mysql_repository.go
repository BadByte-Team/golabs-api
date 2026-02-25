// Package refreshtokeninfra implementa la persistencia de refresh tokens en MySQL/MariaDB.
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

// MySQLRefreshTokenRepository implementa RefreshTokenRepository usando MySQL/MariaDB.
type MySQLRefreshTokenRepository struct {
	db *sql.DB
}

// New crea un nuevo MySQLRefreshTokenRepository con la conexion de base de datos indicada.
func New(db *sql.DB) *MySQLRefreshTokenRepository {
	return &MySQLRefreshTokenRepository{db: db}
}

// HashToken calcula el digest SHA-256 del token crudo y lo retorna como string hexadecimal.
// Esta funcion es el punto unico donde se hashean tokens en todo el modulo.
// Entrada:  raw, valor del token en texto plano.
// Salida:   string hexadecimal de 64 caracteres con el SHA-256 del token.
func HashToken(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

// Save inserta un nuevo refresh token en la tabla refresh_tokens.
// El campo revoked_at no se inserta (queda NULL hasta que sea revocado).
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

// GetByTokenHash busca un refresh token por su hash SHA-256.
// Retorna error si no existe ningun token con ese hash.
func (r *MySQLRefreshTokenRepository) GetByTokenHash(ctx context.Context, hash string) (*refreshtokendomain.RefreshToken, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, user_id, token_hash, expires_at, created_at, revoked_at
		 FROM refresh_tokens WHERE token_hash = ?`,
		hash,
	)
	return scanRefreshToken(row)
}

// Revoke marca el token identificado por id como revocado estableciendo revoked_at = NOW().
// Una vez revocado, el token no puede volver a ser usado.
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

// RevokeAllForUser revoca todos los tokens activos (revoked_at IS NULL) del usuario indicado.
// Se usa al cambiar contrasena o al detectar actividad sospechosa para cerrar todas las sesiones.
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

// scanRefreshToken mapea una fila SQL a un RefreshToken.
// Parsea los UUID de string a uuid.UUID y maneja el campo nullable revoked_at.
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
