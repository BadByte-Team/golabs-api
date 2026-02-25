// Package refreshtokenapp contiene los casos de uso del ciclo de vida de refresh tokens:
// emision, rotacion (refresh) y revocacion (logout).
package refreshtokenapp

import (
	"context"

	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
	refreshtokeninfra "golabs-api/internal/refreshtoken/infrastructure"
)

// RevokeRefreshTokenUseCase invalida un refresh token especifico (logout).
// Su comportamiento es idempotente: si el token no existe o ya fue revocado, no retorna error.
type RevokeRefreshTokenUseCase struct {
	repo refreshtokendomain.RefreshTokenRepository
}

// NewRevokeRefreshTokenUseCase crea un RevokeRefreshTokenUseCase con el repositorio indicado.
func NewRevokeRefreshTokenUseCase(repo refreshtokendomain.RefreshTokenRepository) *RevokeRefreshTokenUseCase {
	return &RevokeRefreshTokenUseCase{repo: repo}
}

// Execute revoca el token que coincide con el valor crudo indicado.
//
// Si el token no se encuentra, se considera que ya estaba revocado y no se retorna error.
// Esto hace que el logout sea siempre exitoso para el cliente, independientemente
// del estado del token (expirado, ya revocado, inexistente).
//
// Entrada:  ctx, contexto de la peticion; rawToken, valor del token en texto plano.
// Salida:   nil en todos los casos salvo error inesperado de base de datos al revocar.
func (uc *RevokeRefreshTokenUseCase) Execute(ctx context.Context, rawToken string) error {
	hash := refreshtokeninfra.HashToken(rawToken)
	rt, err := uc.repo.GetByTokenHash(ctx, hash)
	if err != nil {
		// Token no encontrado: tratar como ya revocado, no retornar error.
		return nil
	}
	return uc.repo.Revoke(ctx, rt.ID)
}
