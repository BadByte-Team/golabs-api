package refreshtokenapp

import (
	"context"

	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
	refreshtokeninfra "golabs-api/internal/refreshtoken/infrastructure"
)

// RevokeRefreshTokenUseCase revokes a single refresh token (logout).
type RevokeRefreshTokenUseCase struct {
	repo refreshtokendomain.RefreshTokenRepository
}

// NewRevokeRefreshTokenUseCase creates a RevokeRefreshTokenUseCase.
func NewRevokeRefreshTokenUseCase(repo refreshtokendomain.RefreshTokenRepository) *RevokeRefreshTokenUseCase {
	return &RevokeRefreshTokenUseCase{repo: repo}
}

// Execute revokes the token matching the given raw value.
// Errors are silently ignored so logout is always a no-op for the client.
func (uc *RevokeRefreshTokenUseCase) Execute(ctx context.Context, rawToken string) error {
	hash := refreshtokeninfra.HashToken(rawToken)
	rt, err := uc.repo.GetByTokenHash(ctx, hash)
	if err != nil {
		// Token not found — treat as already revoked, return no error.
		return nil
	}
	return uc.repo.Revoke(ctx, rt.ID)
}
