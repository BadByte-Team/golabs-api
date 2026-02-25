// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	"golabs-api/internal/infrastructure/security"
	userdomain "golabs-api/internal/user/domain"
)

// LoginUseCase maneja la autenticacion de usuarios con email o username + contrasena.
type LoginUseCase struct {
	repo userdomain.UserRepository
	jwt  *security.JWTService
}

// NewLoginUseCase crea un LoginUseCase con el repositorio y el servicio JWT indicados.
func NewLoginUseCase(repo userdomain.UserRepository, jwt *security.JWTService) *LoginUseCase {
	return &LoginUseCase{repo: repo, jwt: jwt}
}

// Execute autentica al usuario y retorna un access token firmado junto al UUID del usuario.
//
// El campo identifier acepta tanto una direccion de email (si contiene "@")
// como un nombre de usuario. Esto permite que el cliente use un solo campo para ambos.
//
// El UUID del usuario se retorna adicionalmente al token porque el handler lo necesita
// para emitir el refresh token sin una consulta extra a la base de datos.
//
// Seguridad: tanto "usuario no existe" como "contrasena incorrecta" retornan el mismo
// error generico para no revelar si un email/username esta registrado en el sistema.
//
// Entrada:  identifier (email o username), password en texto plano.
// Salida:   accessToken JWT firmado, userID UUID del usuario, o error.
func (uc *LoginUseCase) Execute(identifier, password string) (accessToken string, userID uuid.UUID, err error) {
	var user *userdomain.User

	// Determinar estrategia de busqueda: los emails contienen "@".
	if strings.Contains(identifier, "@") {
		user, err = uc.repo.GetByEmail(identifier)
	} else {
		user, err = uc.repo.GetByUsername(identifier)
	}

	if err != nil {
		// Error generico: no revelar si el email/username existe o no.
		return "", uuid.Nil, errors.New("credenciales invalidas")
	}

	if user.Banned {
		return "", uuid.Nil, errors.New("usuario baneado")
	}

	// Comparacion de hash con bcrypt en tiempo constante para resistir timing attacks.
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(password),
	); err != nil {
		return "", uuid.Nil, errors.New("credenciales invalidas")
	}

	token, err := uc.jwt.Generate(user.ID.String(), user.Role)
	if err != nil {
		return "", uuid.Nil, err
	}

	return token, user.ID, nil
}
