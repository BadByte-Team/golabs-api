package userapp

import (
	userdomain "golabs-api/internal/user/domain"
)

// ListUsersUseCase returns a paginated list of users (admin only).
type ListUsersUseCase struct {
	repo userdomain.UserRepository
}

func NewListUsersUseCase(repo userdomain.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{repo: repo}
}

// Execute returns up to `size` users starting at `page`.
// total is the overall count of users in the system.
func (uc *ListUsersUseCase) Execute(page, size int) ([]*userdomain.User, int, error) {
	offset := (page - 1) * size
	return uc.repo.List(offset, size)
}
