package userhttp

import (
	"encoding/json"
	"net/http"

	"golabs-api/internal/apperrors"
	userapp "golabs-api/internal/user/application"
)

type AuthHandler struct {
	login      *userapp.LoginUseCase
	createUser *userapp.CreateUserUseCase
}

func NewAuthHandler(login *userapp.LoginUseCase, createUser *userapp.CreateUserUseCase) *AuthHandler {
	return &AuthHandler{login: login, createUser: createUser}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := decodeJSON(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	// Identifier can be either an email address or a username.
	token, err := h.login.Execute(req.Identifier, req.Password)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, LoginResponse{Token: token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := decodeJSON(r, &req); err != nil {
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

/* -------- helpers -------- */

// decodeJSON is a shared helper to decode a JSON request body.
func decodeJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}
