package userhttp

import "time"

// Auth
type LoginRequest struct {
	Identifier string `json:"identifier" validate:"required"`
	Password   string `json:"password"   validate:"required,min=6"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int    `json:"expires_in"` // seconds until access_token expires
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" validate:"required"`
}

// User CRUD
type CreateUserRequest struct {
	Username string `json:"username" validate:"required,min=3,max=30"`
	Email    string `json:"email"    validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

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

type UserPublicResponse struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Role     string `json:"role"`
	Points   int    `json:"points"`
}

type UpdateUserRequest struct {
	Username string `json:"username,omitempty" validate:"omitempty,min=3,max=30"`
	Email    string `json:"email,omitempty"   validate:"omitempty,email"`
}

type ChangePasswordRequest struct {
	CurrentPassword string `json:"current_password" validate:"required"`
	NewPassword     string `json:"new_password"     validate:"required,min=6"`
}

type UpdateUserRoleRequest struct {
	Role string `json:"role" validate:"required,oneof=admin user"`
}

type UpdateUserPointsRequest struct {
	Points int `json:"points" validate:"min=0"`
}

type BanUserResponse struct {
	Banned   bool       `json:"banned"`
	BannedAt *time.Time `json:"banned_at,omitempty"`
}
