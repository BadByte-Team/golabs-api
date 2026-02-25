// Package userapp contiene los casos de uso del modulo de usuarios.
package userapp

import (
	"errors"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	userdomain "golabs-api/internal/user/domain"
)

// ChangePasswordUseCase permite a un usuario cambiar su contrasena verificando primero la actual.
// Solo el propio usuario puede cambiar su contrasena; el admin puede hacerlo con UpdatePassword.
type ChangePasswordUseCase struct {
	repo userdomain.UserRepository
}

// NewChangePasswordUseCase crea un ChangePasswordUseCase con el repositorio indicado.
func NewChangePasswordUseCase(repo userdomain.UserRepository) *ChangePasswordUseCase {
	return &ChangePasswordUseCase{repo: repo}
}

// Execute verifica la contrasena actual y la reemplaza con el hash de la nueva.
//
// Seguridad: requiere que el usuario conozca su contrasena actual para poder cambiarla,
// previniendo que un token robado permita tomar control total de la cuenta.
//
// Entrada:  userID (UUID en string), currentPassword y newPassword en texto plano.
// Salida:   error si el ID es invalido, las contrasenas son incorrectas o falla la BD.
func (uc *ChangePasswordUseCase) Execute(userID, currentPassword, newPassword string) error {
	if _, err := uuid.Parse(userID); err != nil {
		return errors.New("id invalido")
	}

	if currentPassword == "" || newPassword == "" {
		return errors.New("contrasenas requeridas")
	}

	if len(newPassword) < 8 {
		return errors.New("la nueva contrasena debe tener al menos 8 caracteres")
	}

	user, err := uc.repo.GetByID(userID)
	if err != nil {
		return err
	}

	// Verificar contrasena actual con bcrypt en tiempo constante.
	if err := bcrypt.CompareHashAndPassword(
		[]byte(user.PasswordHash),
		[]byte(currentPassword),
	); err != nil {
		return errors.New("contrasena actual incorrecta")
	}

	newHash, err := bcrypt.GenerateFromPassword([]byte(newPassword), 12)
	if err != nil {
		return err
	}

	return uc.repo.UpdatePassword(userID, string(newHash))
}
