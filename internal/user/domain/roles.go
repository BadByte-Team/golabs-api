// Package userdomain define el modelo de dominio, las interfaces de repositorio y las constantes del modulo de usuarios.
package userdomain

// Roles de usuario disponibles en el sistema.
// Usar siempre estas constantes en lugar de strings literales
// para evitar errores tipograficos y facilitar el refactoring.
const (
	RoleAdmin = "admin"
	RoleUser  = "user"
)
