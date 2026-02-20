package userhttp

import (
	"encoding/json"
	"net/http"

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
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}
	token, err := h.login.Execute(req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{Token: token})
}

func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	var req CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}
	user, err := h.createUser.Execute(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Points:   user.Points,
	})
}
