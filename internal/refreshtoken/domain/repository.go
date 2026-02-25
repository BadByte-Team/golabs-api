// Package refreshtokendomain define el modelo de dominio del refresh token
// y la interfaz de repositorio para su persistencia.
package refreshtokendomain

import (
	"context"

	"github.com/google/uuid"
)

// RefreshTokenRepository define las operaciones de persistencia de refresh tokens.
// Las implementaciones concretas viven en el paquete infrastructure.
type RefreshTokenRepository interface {
	// Save persiste un nuevo refresh token. El token ya debe tener el hash calculado.
	Save(ctx context.Context, rt *RefreshToken) error

	// GetByTokenHash busca un refresh token por su hash SHA-256.
	// Retorna error si no existe ningun token con ese hash.
	GetByTokenHash(ctx context.Context, hash string) (*RefreshToken, error)

	// Revoke marca un token especifico como revocado, invalidandolo para futuros usos.
	// Se llama tanto al rotar el token (en /auth/refresh) como al hacer logout.
	Revoke(ctx context.Context, id uuid.UUID) error

	// RevokeAllForUser revoca todos los tokens activos de un usuario.
	// Util al cambiar contrasena o al detectar actividad sospechosa.
	RevokeAllForUser(ctx context.Context, userID uuid.UUID) error
}
