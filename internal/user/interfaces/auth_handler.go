package userhttp

import (
	"net/http"

	"golabs-api/internal/apperrors"
	"golabs-api/internal/interfaces/http/validate"
	refreshtokenapp "golabs-api/internal/refreshtoken/application"
	userapp "golabs-api/internal/user/application"
)

// AuthHandler handles public authentication endpoints.
type AuthHandler struct {
	login      *userapp.LoginUseCase
	createUser *userapp.CreateUserUseCase
	issueRT    *refreshtokenapp.IssueRefreshTokenUseCase
	refreshRT  *refreshtokenapp.RefreshAccessTokenUseCase
	revokeRT   *refreshtokenapp.RevokeRefreshTokenUseCase
}

// NewAuthHandler creates an AuthHandler.
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

// Login handles POST /auth/login.
// Returns { access_token, refresh_token, expires_in }.
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
		ExpiresIn:    15 * 60, // 15 minutes in seconds
	})
}

// Register handles POST /auth/register.
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

// Refresh handles POST /auth/refresh.
// Validates the refresh token, rotates it, and returns a new token pair.
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

// Logout handles POST /auth/logout.
// Revokes the refresh token so it can no longer be used.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	var req RefreshRequest
	if err := validate.DecodeAndValidate(r, &req); err != nil {
		// Even with a bad body, respond 204 — logout should never fail for the client.
		w.WriteHeader(http.StatusNoContent)
		return
	}
	_ = h.revokeRT.Execute(r.Context(), req.RefreshToken)
	w.WriteHeader(http.StatusNoContent)
}

/* -------- helpers -------- */

// decodeJSON is a shared helper to decode a JSON request body (no validation).
func decodeJSON(r *http.Request, v any) error {
	return validate.DecodeOnly(r, v)
}
