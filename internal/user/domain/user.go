// Package userdomain define el modelo de dominio del usuario y la interfaz
// de acceso a datos correspondiente.
package userdomain

import (
	"time"

	"github.com/google/uuid"
)

// User representa un usuario registrado en la plataforma.
// El campo PasswordHash nunca se expone en respuestas HTTP; se usa solo internamente
// para la verificacion de credenciales.
type User struct {
	ID           uuid.UUID
	Username     string
	Email        string
	PasswordHash string // hash bcrypt; nunca se serializa hacia el cliente
	Role         string // "admin" | "user"
	Points       int    // puntos acumulados resolviendo challenges

	Banned   bool       // indica si el usuario esta baneado
	BannedAt *time.Time // momento del baneo; nil si no esta baneado

	CreatedAt time.Time
	UpdatedAt time.Time
}
