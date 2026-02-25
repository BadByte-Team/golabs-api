// Package refreshtokendomain define el modelo de dominio del refresh token
// y la interfaz de repositorio para su persistencia.
//
// Los refresh tokens son opacos (secuencias aleatorias) y se almacenan como
// hash SHA-256 en la base de datos. Implementan rotacion automatica: cada uso
// invalida el token anterior y emite uno nuevo.
package refreshtokendomain

import (
	"time"

	"github.com/google/uuid"
)

// RefreshToken representa un refresh token persistido en base de datos.
// El valor en texto plano NUNCA se almacena; solo su hash SHA-256.
//
// Ciclo de vida: emitido al login -> usado en /auth/refresh (queda revocado) ->
// nuevo token emitido -> etc. Expira despues de RefreshTokenTTL dias sin importar uso.
type RefreshToken struct {
	ID        uuid.UUID
	UserID    uuid.UUID
	TokenHash string // digest SHA-256 en hexadecimal del token en texto plano
	ExpiresAt time.Time
	CreatedAt time.Time
	RevokedAt *time.Time // nil si el token sigue vigente; timestamp si fue revocado
}

// IsValid retorna true si el token no fue revocado y no ha expirado.
// Usar este metodo antes de aceptar un refresh token en el use case.
func (rt *RefreshToken) IsValid() bool {
	return rt.RevokedAt == nil && time.Now().Before(rt.ExpiresAt)
}
