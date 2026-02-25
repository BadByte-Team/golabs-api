// Package userhttp implementa los handlers HTTP del modulo de usuarios
// y el registro de sus rutas en el router principal.
package userhttp

import "time"

// ── DTOs de autenticacion ──────────────────────────────────────────────────────

// LoginRequest es el body de la peticion de inicio de sesion.
// El campo identifier acepta tanto email como username.
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password"   validate:"required,min=6"`
}

// LoginResponse contiene el par de tokens emitido al autenticarse exitosamente.
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // segundos hasta que expira el access_token
}

// RefreshRequest es el body de la peticion de rotacion de refresh token.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// ── DTOs de usuario ────────────────────────────────────────────────────────────

// CreateUserRequest es el body para crear un nuevo usuario.
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

// UserResponse es la representacion JSON completa de un usuario (para admins y el propio usuario).
type UserResponse struct {
	ID        string     `json:"id"`
	Username  string     `json:"username"`
	Email     string     `json:"email"`
	Role      string     `json:"role"`
	Points    int        `json:"points"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	Banned    bool       `json:"banned"`
	BannedAt  *time.Time `json:"banned_at,omitempty"`
}

// UserPublicResponse es la representacion JSON reducida de un usuario para respuestas publicas.
// Omite campos sensibles como email y estado de baneo.
type UserPublicResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Points   int    `json:"points"`
}

// UpdateUserRequest es el body para modificar perfil de usuario (patch semantics).
// Los campos vacios u omitidos no se aplican.
type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email    string `json:"email,omitempty"   validate:"omitempty,email"`
}

// ChangePasswordRequest es el body para cambiar la contrasena del usuario.
// Requiere la contrasena actual para confirmar identidad.
type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=6"`
}

// UpdateUserRoleRequest es el body para cambiar el rol de un usuario (admin only).
type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user"`
}

// UpdateUserPointsRequest es el body para establecer los puntos de un usuario (admin only).
type UpdateUserPointsRequest struct {
	Points int `json:"points" validate:"min=0"`
}

// BanUserResponse confirma el estado de baneo del usuario tras una operacion ban/unban.
type BanUserResponse struct {
	Banned   bool       `json:"banned"`
	BannedAt *time.Time `json:"banned_at,omitempty"`
}
