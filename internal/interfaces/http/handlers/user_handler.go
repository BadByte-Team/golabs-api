package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/chi/v5"

	usecases "golabs-api/internal/application/usecases/user"
	"golabs-api/internal/domain/entities"
	"golabs-api/internal/interfaces/dto"
)

type UserHandler struct {
	createUser     *usecases.CreateUserUseCase
	getUserByID    *usecases.GetUserByIDUseCase
	updateUser     *usecases.UpdateUserUseCase
	changePassword *usecases.ChangePasswordUseCase

	updateRole   *usecases.UpdateUserRoleUseCase
	updatePoints *usecases.UpdateUserPointsUseCase
	banUser      *usecases.BanUserUseCase
	unbanUser    *usecases.UnbanUserUseCase
}

func NewUserHandler(
	createUser *usecases.CreateUserUseCase,
	getUserByID *usecases.GetUserByIDUseCase,
	updateUser *usecases.UpdateUserUseCase,
	changePassword *usecases.ChangePasswordUseCase,
	updateRole *usecases.UpdateUserRoleUseCase,
	updatePoints *usecases.UpdateUserPointsUseCase,
	banUser *usecases.BanUserUseCase,
	unbanUser *usecases.UnbanUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUser:     createUser,
		getUserByID:    getUserByID,
		updateUser:     updateUser,
		changePassword: changePassword,
		updateRole:     updateRole,
		updatePoints:   updatePoints,
		banUser:        banUser,
		unbanUser:      unbanUser,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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

	h.respondUser(w, user, http.StatusCreated)
}

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		http.Error(w, "id requerido", http.StatusBadRequest)
		return
	}

	user, err := h.getUserByID.Execute(id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	h.respondUser(w, user, http.StatusOK)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	user, err := h.updateUser.Execute(
		id,
		req.Username,
		req.Email,
	)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	h.respondUser(w, user, http.StatusOK)
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.ChangePasswordRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	if err := h.changePassword.Execute(
		id,
		req.CurrentPassword,
		req.NewPassword,
	); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateUserRoleRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	if err := h.updateRole.Execute(id, req.Role); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdatePoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	var req dto.UpdateUserPointsRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "json inválido", http.StatusBadRequest)
		return
	}

	if err := h.updatePoints.Execute(id, req.Points); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Ban(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.banUser.Execute(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.BanUserResponse{
		Banned: true,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) Unban(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")

	if err := h.unbanUser.Execute(id); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	resp := dto.BanUserResponse{
		Banned: false,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (h *UserHandler) respondUser(
	w http.ResponseWriter,
	user *entities.User,
	status int,
) {
	resp := dto.UserResponse{
		ID:        user.ID.String(),
		Username:  user.Username,
		Email:     user.Email,
		Role:      user.Role,
		Points:    user.Points,
		Banned:    user.Banned,
		BannedAt:  user.BannedAt,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(resp)
}
