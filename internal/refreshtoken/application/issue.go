package refreshtokenapp

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"time"

	"github.com/google/uuid"

	refreshtokeninfra "golabs-api/internal/refreshtoken/infrastructure"

	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
)

const refreshTokenDays = 30

// IssueRefreshTokenUseCase generates and persists a new refresh token for a user.
type IssueRefreshTokenUseCase struct {
	repo refreshtokendomain.RefreshTokenRepository
}

// NewIssueRefreshTokenUseCase creates an IssueRefreshTokenUseCase.
func NewIssueRefreshTokenUseCase(repo refreshtokendomain.RefreshTokenRepository) *IssueRefreshTokenUseCase {
	return &IssueRefreshTokenUseCase{repo: repo}
}

// Execute generates a random raw token, stores its hash, and returns the raw value.
// The caller MUST send the raw token to the client — it is never stored.
func (uc *IssueRefreshTokenUseCase) Execute(ctx context.Context, userID uuid.UUID) (rawToken string, expiresAt time.Time, err error) {
	raw := make([]byte, 32)
	if _, err = rand.Read(raw); err != nil {
		return "", time.Time{}, fmt.Errorf("generate random bytes: %w", err)
	}
	rawToken = hex.EncodeToString(raw)

	expiresAt = time.Now().Add(time.Duration(refreshTokenDays) * 24 * time.Hour)

	rt := &refreshtokendomain.RefreshToken{
		ID:        uuid.New(),
		UserID:    userID,
		TokenHash: refreshtokeninfra.HashToken(rawToken),
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err = uc.repo.Save(ctx, rt); err != nil {
		return "", time.Time{}, err
	}

	return rawToken, expiresAt, nil
}
