package handlers

import (
	"encoding/json"
	"net/http"

	usecases "golabs-api/internal/application/usecases/user"
	"golabs-api/internal/interfaces/dto"
)

type UserHandler struct {
	createUser  *usecases.CreateUserUseCase
	getUserByID *usecases.GetUserByIDUseCase
}

func NewUserHandler(
	createUser *usecases.CreateUserUseCase,
	getUserByID *usecases.GetUserByIDUseCase,
) *UserHandler {
	return &UserHandler{
		createUser:  createUser,
		getUserByID: getUserByID,
	}
}

func (uh *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json invalid", http.StatusBadRequest)
		return
	}

	user, err := uh.createUser.Execute(req.Username, req.Email, req.Password)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Points:   user.Points,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(resp)
}

func (uh *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	id := r.URL.Path[len("/users/"):]
	if id == "" {
		http.Error(w, "id requerido", http.StatusBadRequest)
		return
	}

	user, err := uh.getUserByID.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	resp := dto.UserResponse{
		ID:       user.ID.String(),
		Username: user.Username,
		Email:    user.Email,
		Role:     user.Role,
		Points:   user.Points,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
