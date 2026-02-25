// Package userdomain define el modelo de dominio del usuario y la interfaz
// de acceso a datos correspondiente.
package userdomain

// UserRepository define los metodos de persistencia del agregado User.
// Las implementaciones concretas viven en el paquete infrastructure.
type UserRepository interface {
	// Create persiste un nuevo usuario. Retorna conflicto si email o username ya existen.
	Create(user *User) error

	// GetByID busca un usuario por su UUID en formato string.
	// Retorna error si no existe.
	GetByID(id string) (*User, error)

	// GetByEmail busca un usuario por su direccion de correo.
	GetByEmail(email string) (*User, error)

	// GetByUsername busca un usuario por su nombre de usuario exacto (case-sensitive).
	GetByUsername(username string) (*User, error)

	// SearchByUsername retorna usuarios cuyo username contiene la cadena query (case-insensitive).
	// Util para funciones de busqueda/autocomplete.
	SearchByUsername(query string) ([]*User, error)

	// List retorna una lista paginada de todos los usuarios y el total de registros.
	// offset es el numero de registros a saltar; size es el maximo de resultados.
	List(offset, size int) ([]*User, int, error)

	// Update persiste los cambios de los campos editables de un usuario.
	Update(user *User) error

	// UpdatePassword reemplaza el hash de contrasena de un usuario.
	UpdatePassword(userID string, passwordHash string) error

	// UpdateRole cambia el rol del usuario ("admin" | "user").
	UpdateRole(userID string, role string) error

	// UpdatePoints suma o fija los puntos acumulados del usuario.
	UpdatePoints(userID string, points int) error

	// Ban marca al usuario como baneado; el acceso queda bloqueado inmediatamente.
	Ban(userID string) error

	// Unban revierte el baneo de un usuario.
	Unban(userID string) error
}
