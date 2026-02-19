package db

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	"golabs-api/internal/domain/entities"
	"golabs-api/internal/domain/repositories"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) repositories.UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(user *entities.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password_hash, role, points, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now().UTC()

	user.ID = uuid.New()
	user.CreatedAt = now
	user.UpdatedAt = now

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = smt.Exec(
		user.ID.String(),
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Points,
		user.CreatedAt,
		user.UpdatedAt,
	)

	return err
}

func (r *MySQLUserRepository) GetByEmail(email string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at
		FROM users
		WHERE email = ?
		LIMIT 1
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := smt.QueryRow(email)
	return scanUser(row)
}

func (r *MySQLUserRepository) GetByID(id string) (*entities.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at
		FROM users
		WHERE id = ?
		LIMIT 1
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}

	row := smt.QueryRow(id)
	return scanUser(row)
}

func (r *MySQLUserRepository) Update(user *entities.User) error {
	query := `
		UPDATE users
		SET username = ?, email = ?, password_hash = ?, role = ?, points = ?, updated_at = ?
		WHERE id = ?
	`

	now := time.Now().UTC()

	user.UpdatedAt = now

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = smt.Exec(
		user.Username,
		user.Email,
		user.PasswordHash,
		user.Role,
		user.Points,
		user.UpdatedAt,
		user.ID.String(),
	)

	return err
}

func (r *MySQLUserRepository) Delete(id string) error {
	query := `
		DELETE FROM users
		WHERE id = ?
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}

	_, err = smt.Exec(id)
	return err
}

func scanUser(row *sql.Row) (*entities.User, error) {
	var user entities.User
	var id string

	err := row.Scan(
		&id,
		&user.Username,
		&user.Email,
		&user.PasswordHash,
		&user.Role,
		&user.Points,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	user.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	return &user, nil
}
