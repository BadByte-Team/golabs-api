// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	userdomain "golabs-api/internal/user/domain"
)

// ListUsersUseCase retorna una lista paginada de todos los usuarios del sistema.
// Endpoint de solo administrador; los usuarios regulares no pueden listar a otros usuarios.
type ListUsersUseCase struct {
	repo userdomain.UserRepository
}

// NewListUsersUseCase crea un ListUsersUseCase con el repositorio indicado.
func NewListUsersUseCase(repo userdomain.UserRepository) *ListUsersUseCase {
	return &ListUsersUseCase{repo: repo}
}

// Execute retorna la pagina indicada de usuarios y el total de registros en el sistema.
//
// Entrada:  page (base 1), size (cantidad de registros por pagina).
// Salida:   slice de usuarios, total de registros para calcular paginas, o error.
func (uc *ListUsersUseCase) Execute(page, size int) ([]*userdomain.User, int, error) {
	offset := (page - 1) * size
	return uc.repo.List(offset, size)
}
