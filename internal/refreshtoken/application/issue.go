// Package refreshtokenapp contiene los casos de uso del ciclo de vida de refresh tokens:
// emision, rotacion (refresh) y revocacion (logout).
package refreshtokenapp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
	refreshtokeninfra "golabs-api/internal/refreshtoken/infrastructure"
)

// refreshTokenDays es el tiempo de vida del refresh token en dias.
// Al expirar, el usuario debe volver a autenticarse con su contrasena.
const refreshTokenDays = 30

// IssueRefreshTokenUseCase genera y persiste un nuevo refresh token para un usuario.
// Se llama inmediatamente despues del login exitoso.
type IssueRefreshTokenUseCase struct {
	repo refreshtokendomain.RefreshTokenRepository
}

// NewIssueRefreshTokenUseCase crea un IssueRefreshTokenUseCase con el repositorio indicado.
func NewIssueRefreshTokenUseCase(repo refreshtokendomain.RefreshTokenRepository) *IssueRefreshTokenUseCase {
	return &IssueRefreshTokenUseCase{repo: repo}
}

// Execute genera un token opaco aleatorio, almacena su hash SHA-256 y retorna el valor crudo.
//
// El valor crudo (rawToken) se envia al cliente UNA SOLA VEZ y nunca se almacena en BD.
// Solo el hash se persiste para poder validarlo en futuros requests sin exponer el token original.
//
// Entrada:  ctx, contexto de la peticion; userID, UUID del usuario autenticado.
// Salida:   rawToken (64 chars hex), expiresAt, o error si falla la generacion o la BD.
func (uc *IssueRefreshTokenUseCase) Execute(ctx context.Context, userID uuid.UUID) (rawToken string, expiresAt time.Time, err error) {
	// Generar 32 bytes aleatorios criptograficamente seguros.
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", time.Time{}, fmt.Errorf("generate random bytes: %w", err)
	}
	rawToken = hex.EncodeToString(raw)

	expiresAt = time.Now().Add(time.Duration(refreshTokenDays) * 24 * time.Hour)

	rt := &refreshtokendomain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: refreshtokeninfra.HashToken(rawToken), // solo el hash va a BD
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err = uc.repo.Save(ctx, rt); err != nil {
		return "", time.Time{}, err
	}

	return rawToken, expiresAt, nil
}
