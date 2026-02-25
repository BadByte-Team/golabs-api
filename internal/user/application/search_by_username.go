// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	userdomain "golabs-api/internal/user/domain"
)

// SearchUserByUsernameUseCase realiza una busqueda parcial de usuarios por nombre de usuario.
// Util para implementar autocomplete o busqueda de jugadores en la UI.
type SearchUserByUsernameUseCase struct {
	repo userdomain.UserRepository
}

// NewSearchUserByUsernameUseCase crea un SearchUserByUsernameUseCase con el repositorio indicado.
func NewSearchUserByUsernameUseCase(repo userdomain.UserRepository) *SearchUserByUsernameUseCase {
	return &SearchUserByUsernameUseCase{repo: repo}
}

// Execute busca usuarios cuyo username contenga el termino de busqueda (LIKE %query%).
// El termino de busqueda es obligatorio para evitar retornar todos los usuarios.
//
// Entrada:  query, termino de busqueda (minimo 1 caracter).
// Salida:   slice de usuarios que coinciden con el termino o error.
func (uc *SearchUserByUsernameUseCase) Execute(query string) ([]*userdomain.User, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}
	return uc.repo.SearchByUsername(query)
}
