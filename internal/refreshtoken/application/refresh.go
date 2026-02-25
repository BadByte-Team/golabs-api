// Package refreshtokenapp contiene los casos de uso del ciclo de vida de refresh tokens:
// emision, rotacion (refresh) y revocacion (logout).
package refreshtokenapp

import (
	"context"
	"errors"

	refreshtokeninfra "golabs-api/internal/refreshtoken/infrastructure"

	"golabs-api/internal/infrastructure/security"
	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
	userdomain "golabs-api/internal/user/domain"
)

// RefreshAccessTokenUseCase valida un refresh token, lo rota y emite un nuevo par de tokens.
//
// Token rotation: cada uso del refresh token lo invalida y genera uno nuevo.
// Si un refresh token fue robado y el atacante intenta usarlo despues de que el usuario
// legitimo lo uso, la validacion fallara (el token ya fue revocado).
type RefreshAccessTokenUseCase struct {
	rtRepo   refreshtokendomain.RefreshTokenRepository
	userRepo userdomain.UserRepository
	jwt      *security.JWTService
	issue    *IssueRefreshTokenUseCase // para emitir el nuevo refresh token
}

// NewRefreshAccessTokenUseCase crea un RefreshAccessTokenUseCase con las dependencias indicadas.
func NewRefreshAccessTokenUseCase(
	rtRepo refreshtokendomain.RefreshTokenRepository,
	userRepo userdomain.UserRepository,
	jwt *security.JWTService,
	issue *IssueRefreshTokenUseCase,
) *RefreshAccessTokenUseCase {
	return &RefreshAccessTokenUseCase{rtRepo: rtRepo, userRepo: userRepo, jwt: jwt, issue: issue}
}

// RefreshResult contiene el nuevo par de tokens emitido tras una rotacion exitosa.
type RefreshResult struct {
	AccessToken  string
	RefreshToken string
}

// ErrInvalidRefreshToken se retorna cuando el refresh token no existe, ya fue revocado o expiro.
// Se usa un error generico unico para no revelar si el token existio alguna vez.
var ErrInvalidRefreshToken = errors.New("refresh token invalido o expirado")

// Execute valida el refresh token crudo, lo revoca y emite un nuevo par de tokens.
//
// Flujo:
//  1. Calcular hash SHA-256 del token y buscar en BD.
//  2. Verificar que el token sea valido (no expirado, no revocado).
//  3. Revocar el token actual (token rotation).
//  4. Cargar datos frescos del usuario (rol actualizado, estado de baneo).
//  5. Generar nuevo access token (JWT).
//  6. Emitir nuevo refresh token y persistirlo.
//
// Entrada:  ctx, contexto de la peticion; rawToken, valor del refresh token en texto plano.
// Salida:   RefreshResult con los nuevos tokens, o error si el token es invalido.
func (uc *RefreshAccessTokenUseCase) Execute(ctx context.Context, rawToken string) (*RefreshResult, error) {
	hash := refreshtokeninfra.HashToken(rawToken)

	rt, err := uc.rtRepo.GetByTokenHash(ctx, hash)
	if err != nil || !rt.IsValid() {
		return nil, ErrInvalidRefreshToken
	}

	// Revocar el token actual inmediatamente (token rotation).
	if err := uc.rtRepo.Revoke(ctx, rt.ID); err != nil {
		return nil, err
	}

	// Cargar datos frescos del usuario: el rol pudo haber cambiado desde el login.
	user, err := uc.userRepo.GetByID(rt.UserID.String())
	if err != nil {
		return nil, ErrInvalidRefreshToken
	}
	if user.Banned {
		return nil, errors.New("usuario baneado")
	}

	// Generar nuevo access token con datos actuales del usuario.
	accessToken, err := uc.jwt.Generate(user.ID.String(), user.Role)
	if err != nil {
		return nil, err
	}

	// Emitir nuevo refresh token (persistido en BD).
	newRaw, _, err := uc.issue.Execute(ctx, user.ID)
	if err != nil {
		return nil, err
	}

	return &RefreshResult{AccessToken: accessToken, RefreshToken: newRaw}, nil
}
