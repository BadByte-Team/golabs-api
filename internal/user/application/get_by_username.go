package userapp

import (
	userdomain "golabs-api/internal/user/domain"
)

// GetUserByUsernameUseCase retrieves a single user by exact username.
type GetUserByUsernameUseCase struct {
	repo userdomain.UserRepository
}

func NewGetUserByUsernameUseCase(repo userdomain.UserRepository) *GetUserByUsernameUseCase {
	return &GetUserByUsernameUseCase{repo: repo}
}

func (uc *GetUserByUsernameUseCase) Execute(username string) (*userdomain.User, error) {
	return uc.repo.GetByUsername(username)
}
