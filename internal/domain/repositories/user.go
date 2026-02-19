package repositories

import "golabs-api/internal/domain/entities"

type UserRepository interface {
	Create(user *entities.User) error
	GetByEmail(email string) (*entities.User, error)
	GetByID(id string) (*entities.User, error)
	Update(user *entities.User) error
	Delete(id string) error
}
