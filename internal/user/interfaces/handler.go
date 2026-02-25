package userhttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"golabs-api/internal/apperrors"
	userapp "golabs-api/internal/user/application"
	userdomain "golabs-api/internal/user/domain"
)

type UserHandler struct {
	createUser        *userapp.CreateUserUseCase
	getUserByID       *userapp.GetUserByIDUseCase
	getUserByUsername *userapp.GetUserByUsernameUseCase
	searchByUsername  *userapp.SearchUserByUsernameUseCase
	updateUser        *userapp.UpdateUserUseCase
	changePassword    *userapp.ChangePasswordUseCase
	updateRole        *userapp.UpdateUserRoleUseCase
	updatePoints      *userapp.UpdateUserPointsUseCase
	banUser           *userapp.BanUserUseCase
	unbanUser         *userapp.UnbanUserUseCase
}

func NewUserHandler(
	createUser *userapp.CreateUserUseCase,
	getUserByID *userapp.GetUserByIDUseCase,
	getUserByUsername *userapp.GetUserByUsernameUseCase,
	searchByUsername *userapp.SearchUserByUsernameUseCase,
	updateUser *userapp.UpdateUserUseCase,
	changePassword *userapp.ChangePasswordUseCase,
	updateRole *userapp.UpdateUserRoleUseCase,
	updatePoints *userapp.UpdateUserPointsUseCase,
	banUser *userapp.BanUserUseCase,
	unbanUser *userapp.UnbanUserUseCase,
) *UserHandler {
	return &UserHandler{
		createUser:        createUser,
		getUserByID:       getUserByID,
		getUserByUsername: getUserByUsername,
		searchByUsername:  searchByUsername,
		updateUser:        updateUser,
		changePassword:    changePassword,
		updateRole:        updateRole,
		updatePoints:      updatePoints,
		banUser:           banUser,
		unbanUser:         unbanUser,
	}
}

func (h *UserHandler) Create(w http.ResponseWriter, r *http.Request) {
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

func (h *UserHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if id == "" {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	user, err := h.getUserByID.Execute(id)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapUser(user))
}

// GetByUsername handles GET /users/by-username/{username}
func (h *UserHandler) GetByUsername(w http.ResponseWriter, r *http.Request) {
	username := chi.URLParam(r, "username")
	if username == "" {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	user, err := h.getUserByUsername.Execute(username)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapUser(user))
}

// Search handles GET /users/search?q=<query>
func (h *UserHandler) Search(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query().Get("q")
	if q == "" {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	users, err := h.searchByUsername.Execute(q)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, mapUser(u))
	}
	apperrors.RespondJSON(w, http.StatusOK, resp)
}

func (h *UserHandler) Update(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateUserRequest
	if err := decodeJSON(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	user, err := h.updateUser.Execute(id, req.Username, req.Email)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, mapUser(user))
}

func (h *UserHandler) ChangePassword(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req ChangePasswordRequest
	if err := decodeJSON(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	if err := h.changePassword.Execute(id, req.CurrentPassword, req.NewPassword); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdateRole(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateUserRoleRequest
	if err := decodeJSON(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	if err := h.updateRole.Execute(id, req.Role); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) UpdatePoints(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	var req UpdateUserPointsRequest
	if err := decodeJSON(r, &req); err != nil {
		apperrors.RespondError(w, apperrors.ErrBadRequest)
		return
	}
	if err := h.updatePoints.Execute(id, req.Points); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *UserHandler) Ban(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.banUser.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, BanUserResponse{Banned: true})
}

func (h *UserHandler) Unban(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.unbanUser.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, BanUserResponse{Banned: false})
}

/* -------- helpers -------- */

func mapUser(u *userdomain.User) UserResponse {
	return UserResponse{
		ID:        u.ID.String(),
		Username:  u.Username,
		Email:     u.Email,
		Role:      u.Role,
		Points:    u.Points,
		Banned:    u.Banned,
		BannedAt:  u.BannedAt,
		CreatedAt: u.CreatedAt,
		UpdatedAt: u.UpdatedAt,
	}
}
