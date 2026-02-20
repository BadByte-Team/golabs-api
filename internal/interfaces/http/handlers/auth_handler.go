package handlers

import (
	"encoding/json"
	"net/http"

	usecases "golabs-api/internal/application/usecases/user"
	"golabs-api/internal/interfaces/dto"
)

type AuthHandler struct {
	login      *usecases.LoginUseCase
	createUser *usecases.CreateUserUseCase
}

func NewAuthHandler(login *usecases.LoginUseCase, createUser *usecases.CreateUserUseCase) *AuthHandler {
	return &AuthHandler{login: login, createUser: createUser}
}

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	token, err := h.login.Execute(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}

	resp := dto.LoginResponse{Token: token}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	user, err := h.createUser.Execute(
		req.Username,
		req.Email,
		req.Password,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Points:   user.Points,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}
