package userapp

import (
	"errors"

	userdomain "golabs-api/internal/user/domain"
)

// SearchUserByUsernameUseCase performs a partial-match username search.
type SearchUserByUsernameUseCase struct {
	repo userdomain.UserRepository
}

func NewSearchUserByUsernameUseCase(repo userdomain.UserRepository) *SearchUserByUsernameUseCase {
	return &SearchUserByUsernameUseCase{repo: repo}
}

func (uc *SearchUserByUsernameUseCase) Execute(query string) ([]*userdomain.User, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}
	return uc.repo.SearchByUsername(query)
}
