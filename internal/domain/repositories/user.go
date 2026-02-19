package repositories

import "golabs-api/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	GetByID(id string) (*entities.User, error)
	GetByEmail(email string) (*entities.User, error)

	Update(user *entities.User) error
	UpdatePassword(userID string, passwordHash string) error

	UpdateRole(userID string, role string) error
	UpdatePoints(userID string, points int) error

	Ban(userID string) error
	Unban(userID string) error
}
