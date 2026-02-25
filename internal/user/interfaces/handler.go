// Package userhttp implementa los handlers HTTP del modulo de usuarios
// y el registro de sus rutas en el router principal.
package userhttp

import (
	"net/http"

	"github.com/go-chi/chi/v5"

	"golabs-api/internal/apperrors"
	"golabs-api/internal/interfaces/http/pagination"
	userapp "golabs-api/internal/user/application"
	userdomain "golabs-api/internal/user/domain"
)

// UserHandler agrupa los handlers HTTP del modulo de usuarios (autenticados).
// Los endpoints publicos de autenticacion se manejan en AuthHandler.
type UserHandler struct {
	createUser        *userapp.CreateUserUseCase
	getUserByID       *userapp.GetUserByIDUseCase
	getUserByUsername *userapp.GetUserByUsernameUseCase
	searchByUsername  *userapp.SearchUserByUsernameUseCase
	listUsers         *userapp.ListUsersUseCase
	updateUser        *userapp.UpdateUserUseCase
	changePassword    *userapp.ChangePasswordUseCase
	updateRole        *userapp.UpdateUserRoleUseCase
	updatePoints      *userapp.UpdateUserPointsUseCase
	banUser           *userapp.BanUserUseCase
	unbanUser         *userapp.UnbanUserUseCase
}

// NewUserHandler inyecta todas las dependencias del UserHandler.
func NewUserHandler(
	createUser *userapp.CreateUserUseCase,
	getUserByID *userapp.GetUserByIDUseCase,
	getUserByUsername *userapp.GetUserByUsernameUseCase,
	searchByUsername *userapp.SearchUserByUsernameUseCase,
	listUsers *userapp.ListUsersUseCase,
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
		listUsers:         listUsers,
		updateUser:        updateUser,
		changePassword:    changePassword,
		updateRole:        updateRole,
		updatePoints:      updatePoints,
		banUser:           banUser,
		unbanUser:         unbanUser,
	}
}

// List godoc — GET /api/v1/users (admin)
// Retorna la lista paginada de todos los usuarios del sistema.
// Exito: 200 pagination.Response[UserResponse]
func (h *UserHandler) List(w http.ResponseWriter, r *http.Request) {
	pg := pagination.Parse(r)
	users, total, err := h.listUsers.Execute(pg.Number, pg.Size)
	if err != nil {
		apperrors.RespondError(w, err)
		return
	}

	resp := make([]UserResponse, 0, len(users))
	for _, u := range users {
		resp = append(resp, mapUser(u))
	}
	apperrors.RespondJSON(w, http.StatusOK, pagination.New(resp, pg, total))
}

// Create godoc — POST /api/v1/users (admin)
// Crea un usuario directamente (sin verificacion de email). Igual que /auth/register pero para admins.
// Body: CreateUserRequest | Exito: 201 UserResponse
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

// GetByID godoc — GET /api/v1/users/{id}
// Retorna el perfil de un usuario por su UUID.
// Exito: 200 UserResponse | Error: 404 si no existe
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

// GetByUsername godoc — GET /api/v1/users/by-username/{username}
// Retorna el perfil de un usuario por su username exacto.
// Exito: 200 UserResponse | Error: 404 si no existe
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

// Search godoc — GET /api/v1/users/search?q=<query>
// Busca usuarios cuyo username contenga el termino de busqueda (coincidencia parcial).
// Exito: 200 []UserResponse
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

// Update godoc — PATCH /api/v1/users/{id}
// Actualiza username y/o email del usuario (patch semantics).
// Body: UpdateUserRequest | Exito: 200 UserResponse
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

// ChangePassword godoc — PUT /api/v1/users/{id}/password
// Cambia la contrasena del usuario verificando la actual primero.
// Body: ChangePasswordRequest | Exito: 204 No Content
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

// UpdateRole godoc — PUT /api/v1/users/{id}/role (admin)
// Cambia el rol del usuario entre "user" y "admin".
// Body: { "role": "admin"|"user" } | Exito: 204 No Content
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

// UpdatePoints godoc — PUT /api/v1/users/{id}/points (admin)
// Establece los puntos del usuario directamente.
// Body: { "points": int } | Exito: 204 No Content
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

// Ban godoc — POST /api/v1/users/{id}/ban (admin)
// Suspende el acceso del usuario.
// Exito: 200 { "banned": true }
func (h *UserHandler) Ban(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.banUser.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, BanUserResponse{Banned: true})
}

// Unban godoc — POST /api/v1/users/{id}/unban (admin)
// Reactiva el acceso de un usuario baneado.
// Exito: 200 { "banned": false }
func (h *UserHandler) Unban(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	if err := h.unbanUser.Execute(id); err != nil {
		apperrors.RespondError(w, err)
		return
	}
	apperrors.RespondJSON(w, http.StatusOK, BanUserResponse{Banned: false})
}

// mapUser convierte un User de dominio a su representacion JSON para la API.
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
