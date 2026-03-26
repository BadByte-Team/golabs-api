// Package infrastructure implementa el repositorio de eventos usando MySQL/MariaDB.
package infrastructure

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

// MySQLEventRepository implementa eventdomain.Repository usando MySQL/MariaDB.
type MySQLEventRepository struct {
	db *sql.DB
}

// NewEventRepository crea una instancia de MySQLEventRepository a partir de la conexion de BD.
// Retorna la interfaz eventdomain.Repository para que el llamador no dependa de la implementacion concreta.
func NewEventRepository(db *sql.DB) eventdomain.Repository {
	return &MySQLEventRepository{db: db}
}

// Save inserta un nuevo evento en la tabla events.
// Establece los timestamps created_at y updated_at en UTC al momento de la insercion.
func (r *MySQLEventRepository) Save(event *eventdomain.Event) error {
	query := `
		INSERT INTO events (
			id, name, description, max_team_size, status,
			starts_at, ends_at, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now().UTC()
	event.CreatedAt = now
	event.UpdatedAt = now

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		event.ID.String(),
		event.Name,
		event.Description,
		event.MaxTeamSize,
		string(event.Status),
		event.StartsAt,
		event.EndsAt,
		event.CreatedAt,
		event.UpdatedAt,
	)

	return err
}

// GetByID busca un evento por su UUID.
// Retorna error si no existe ningun evento con ese ID.
func (r *MySQLEventRepository) GetByID(id uuid.UUID) (*eventdomain.Event, error) {
	query := `
		SELECT id, name, description, max_team_size, status,
		       starts_at, ends_at, created_at, updated_at
		FROM events
		WHERE id = ?
		LIMIT 1
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanEvent(smt.QueryRow(id.String()))
}

// List retorna todos los eventos ordenados por fecha de inicio descendente.
func (r *MySQLEventRepository) List() ([]*eventdomain.Event, error) {
	query := `
		SELECT id, name, description, max_team_size, status,
		       starts_at, ends_at, created_at, updated_at
		FROM events
		ORDER BY starts_at DESC
	`

	rows, err := r.db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*eventdomain.Event
	for rows.Next() {
		event, err := scanEventRow(rows)
		if err != nil {
			return nil, err
		}
		events = append(events, event)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return events, nil
}

// Update persiste los campos modificables de un evento existente.
// Actualiza automaticamente el campo updated_at al momento de la operacion.
func (r *MySQLEventRepository) Update(event *eventdomain.Event) error {
	query := `
		UPDATE events
		SET name = ?, description = ?, max_team_size = ?, status = ?,
		    starts_at = ?, ends_at = ?, updated_at = ?
		WHERE id = ?
	`

	event.UpdatedAt = time.Now().UTC()

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		event.Name,
		event.Description,
		event.MaxTeamSize,
		string(event.Status),
		event.StartsAt,
		event.EndsAt,
		event.UpdatedAt,
		event.ID.String(),
	)

	return err
}

// scanEvent mapea un sql.Row a un Event, convirtiendo el UUID y el status de string.
func scanEvent(row *sql.Row) (*eventdomain.Event, error) {
	var e eventdomain.Event
	var id string
	var status string

	err := row.Scan(
		&id,
		&e.Name,
		&e.Description,
		&e.MaxTeamSize,
		&status,
		&e.StartsAt,
		&e.EndsAt,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("event not found")
	}
	if err != nil {
		return nil, err
	}

	e.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	e.Status = eventdomain.EventStatus(status)
	return &e, nil
}

// scanEventRow mapea una fila de sql.Rows a un Event para uso en iteraciones.
func scanEventRow(rows *sql.Rows) (*eventdomain.Event, error) {
	var e eventdomain.Event
	var id string
	var status string

	err := rows.Scan(
		&id,
		&e.Name,
		&e.Description,
		&e.MaxTeamSize,
		&status,
		&e.StartsAt,
		&e.EndsAt,
		&e.CreatedAt,
		&e.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	e.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	e.Status = eventdomain.EventStatus(status)
	return &e, nil
}

// Delete elimina un evento por su UUID.
func (r *MySQLEventRepository) Delete(id uuid.UUID) error {
	query := `DELETE FROM events WHERE id = ?`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(id.String())
	return err
}
