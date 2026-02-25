// Package userhttp implementa los handlers HTTP del modulo de usuarios
// y el registro de sus rutas en el router principal.
package userhttp

import (
	"net/http"

	"golabs-api/internal/apperrors"
	"golabs-api/internal/interfaces/http/validate"
	refreshtokenapp "golabs-api/internal/refreshtoken/application"
	userapp "golabs-api/internal/user/application"
)

// AuthHandler agrupa los handlers de autenticacion publica (sin token requerido).
// Depende de los use cases de login, registro y gestion de refresh tokens.
type AuthHandler struct {
	login      *userapp.LoginUseCase
	createUser *userapp.CreateUserUseCase
	issueRT    *refreshtokenapp.IssueRefreshTokenUseCase  // emitir refresh token al login
	refreshRT  *refreshtokenapp.RefreshAccessTokenUseCase // rotar refresh token
	revokeRT   *refreshtokenapp.RevokeRefreshTokenUseCase // revocar refresh token al logout
}

// NewAuthHandler inyecta las dependencias del AuthHandler.
func NewAuthHandler(
	login *userapp.LoginUseCase,
	createUser *userapp.CreateUserUseCase,
	issueRT *refreshtokenapp.IssueRefreshTokenUseCase,
	refreshRT *refreshtokenapp.RefreshAccessTokenUseCase,
	revokeRT *refreshtokenapp.RevokeRefreshTokenUseCase,
) *AuthHandler {
	return &AuthHandler{
		login:      login,
		createUser: createUser,
		issueRT:    issueRT,
		refreshRT:  refreshRT,
		revokeRT:   revokeRT,
	}
}

// Login godoc
//
// POST /auth/login
//
// Autentica al usuario con identifier (email o username) + password.
// Retorna un par de tokens: access token (JWT de 15 min) y refresh token (opaco, de larga duracion).
//
// Body:   { "identifier": string, "password": string }
// Exito:  200 { "access_token", "refresh_token", "expires_in" }
// Error:  400 (body invalido), 401/403 (credenciales incorrectas o usuario baneado)
func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	accessToken, userID, err := h.login.Execute(req.Identifier, req.Password)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	rawRT, _, err := h.issueRT.Execute(r.Context(), userID)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, LoginResponse{
		AccessToken:  accessToken,
		RefreshToken: rawRT,
		ExpiresIn:    15 * 60, // segundos hasta que expira el access token
	})
}

// Register godoc
//
// POST /auth/register
//
// Crea una cuenta nueva con username, email y password.
//
// Body:   { "username", "email", "password" }
// Exito:  201 UserResponse
// Error:  400 (body invalido), 409 (email o username ya en uso)
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	user, err := h.createUser.Execute(req.Username, req.Email, req.Password)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusCreated, mapUser(user))
}

// Refresh godoc
//
// POST /auth/refresh
//
// Valida el refresh token, lo revoca (token rotation) y emite un nuevo par de tokens.
// Si el refresh token fue robado y ya fue rotado, la validacion fallara.
//
// Body:   { "refresh_token": string }
// Exito:  200 { "access_token", "refresh_token", "expires_in" }
// Error:  400 (body invalido), 401 (token invalido, expirado o ya revocado)
func (h *AuthHandler) Refresh(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}

	result, err := h.refreshRT.Execute(r.Context(), req.RefreshToken)
	if err != nil {
		apperrors.RespondJSON(w, http.StatusUnauthorized, map[string]string{"error": err.Error()})
		return
	}

	apperrors.RespondJSON(w, http.StatusOK, LoginResponse{
		AccessToken:  result.AccessToken,
		RefreshToken: result.RefreshToken,
		ExpiresIn:    15 * 60,
	})
}

// Logout godoc
//
// POST /auth/logout
//
// Revoca el refresh token indicado para invalidar la sesion actual.
// Siempre responde con 204 independientemente del resultado, para que el cliente
// no distinga entre "token valido revocado" y "token ya expirado/invalido".
//
// Body:   { "refresh_token": string }
// Exito:  204 No Content (idempotente)
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_ = h.revokeRT.Execute(r.Context(), req.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

// decodeJSON es un helper interno para decodificar el body JSON sin validacion de struct tags.
func decodeJSON(r *http.Request, v any) error {
	return validate.DecodeOnly(r, v)
}
