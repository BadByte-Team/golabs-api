package domain

import "github.com/google/uuid"

type Repository interface {
	Save(event *Event) error
	GetByID(id uuid.UUID) (*Event, error)
	List() ([]*Event, error)
	Update(event *Event) error
}
