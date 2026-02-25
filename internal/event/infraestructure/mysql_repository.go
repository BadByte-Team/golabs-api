package infrastructure

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	eventdomain "golabs-api/internal/event/domain"
)

type MySQLEventRepository struct {
	db *sql.DB
}

func NewEventRepository(db *sql.DB) eventdomain.Repository {
	return &MySQLEventRepository{db: db}
}

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

/* ---------- helpers ---------- */

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
