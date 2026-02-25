package refreshtokenapp

import (
	"context"
	"errors"

	refreshtokeninfra "golabs-api/internal/refreshtoken/infrastructure"

	"golabs-api/internal/infrastructure/security"
	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
	userdomain "golabs-api/internal/user/domain"
)

// RefreshAccessTokenUseCase validates a refresh token, rotates it, and issues a new access + refresh token pair.
type RefreshAccessTokenUseCase struct {
	rtRepo   refreshtokendomain.RefreshTokenRepository
	userRepo userdomain.UserRepository
	jwt      *security.JWTService
	issue    *IssueRefreshTokenUseCase
}

// NewRefreshAccessTokenUseCase creates a RefreshAccessTokenUseCase.
func NewRefreshAccessTokenUseCase(
	rtRepo refreshtokendomain.RefreshTokenRepository,
	userRepo userdomain.UserRepository,
	jwt *security.JWTService,
	issue *IssueRefreshTokenUseCase,
) *RefreshAccessTokenUseCase {
	return &RefreshAccessTokenUseCase{rtRepo: rtRepo, userRepo: userRepo, jwt: jwt, issue: issue}
}

// RefreshResult holds the new token pair returned after a successful refresh.
type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

// ErrInvalidRefreshToken is returned for any invalid, expired, or revoked refresh token.
var ErrInvalidRefreshToken = errors.New("refresh token inválido o expirado")

// Execute validates the raw refresh token, revokes it (rotation), and issues a new pair.
func (uc *RefreshAccessTokenUseCase) Execute(ctx context.Context, rawToken string) (*RefreshResult, error) {
	hash := refreshtokeninfra.HashToken(rawToken)

	rt, err := uc.rtRepo.GetByTokenHash(ctx, hash)
	if err != nil || !rt.IsValid() {
		return nil, ErrInvalidRefreshToken
	}

	// Revoke the old token immediately (token rotation).
	if err := uc.rtRepo.Revoke(ctx, rt.ID); err != nil {
		return nil, err
	}

	// Load fresh user data to get current role and ban status.
	user, err := uc.userRepo.GetByID(rt.UserID.String())
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	if user.Banned {
		return nil, errors.New("usuario baneado")
	}

	// Issue new access token.
	accessToken, err := uc.jwt.Generate(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	// Issue new refresh token.
	newRaw, _, err := uc.issue.Execute(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{AccessToken: accessToken, RefreshToken: newRaw}, nil
}
