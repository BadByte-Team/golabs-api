// Package domain define los tipos de dominio del modulo de eventos CTF,
// incluyendo el modelo Event, sus estados posibles y las reglas de negocio
// que controlan las transiciones de estado.
package domain

import "github.com/google/uuid"

// Repository define los metodos de persistencia del agregado Event.
// Las implementaciones concretas viven en el paquete infrastructure.
type Repository interface {
	// Save persiste un nuevo evento. Retorna error si ya existe un evento con el mismo ID.
	Save(event *Event) error

	// GetByID busca un evento por su UUID. Retorna error si no existe.
	GetByID(id uuid.UUID) (*Event, error)

	// List retorna todos los eventos, sin paginacion.
	// TODO: agregar paginacion cuando el volumen de eventos lo justifique.
	List() ([]*Event, error)

	// Update persiste los cambios en un evento existente (estado, fechas, etc.).
	Update(event *Event) error

	// Delete elimina un evento por su UUID.
	Delete(id uuid.UUID) error
}
