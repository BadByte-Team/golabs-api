package userinfra

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

type MySQLUserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) userdomain.UserRepository {
	return &MySQLUserRepository{db: db}
}

func (r *MySQLUserRepository) Create(user *userdomain.User) error {
	query := `
		INSERT INTO users (
			id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
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
		user.ID.String(), user.Username, user.Email,
		user.PasswordHash, user.Role, user.Points,
		user.CreatedAt, user.UpdatedAt, user.Banned, user.BannedAt,
	)
	return err
}

func (r *MySQLUserRepository) GetByEmail(email string) (*userdomain.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		FROM users WHERE email = ? LIMIT 1
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return scanUser(smt.QueryRow(email))
}

func (r *MySQLUserRepository) GetByID(id string) (*userdomain.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		FROM users WHERE id = ? LIMIT 1
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	return scanUser(smt.QueryRow(id))
}

func (r *MySQLUserRepository) Update(user *userdomain.User) error {
	query := `UPDATE users SET username = ?, email = ?, points = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	user.UpdatedAt = now
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = smt.Exec(user.Username, user.Email, user.Points, now, user.ID.String())
	return err
}

func (r *MySQLUserRepository) Delete(id string) error {
	query := `DELETE FROM users WHERE id = ?`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = smt.Exec(id)
	return err
}

func (r *MySQLUserRepository) UpdatePassword(userID, passwordHash string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = smt.Exec(passwordHash, now, userID)
	return err
}

func (r *MySQLUserRepository) UpdateRole(userID, role string) error {
	query := `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = smt.Exec(role, now, userID)
	return err
}

func (r *MySQLUserRepository) UpdatePoints(userID string, points int) error {
	query := `UPDATE users SET points = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = smt.Exec(points, now, userID)
	return err
}

func (r *MySQLUserRepository) Ban(userID string) error {
	query := `UPDATE users SET banned = ?, banned_at = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = smt.Exec(true, now, now, userID)
	return err
}

func (r *MySQLUserRepository) Unban(userID string) error {
	query := `UPDATE users SET banned = ?, banned_at = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	_, err = smt.Exec(false, nil, now, userID)
	return err
}

func scanUser(row *sql.Row) (*userdomain.User, error) {
	var u userdomain.User
	var id string
	var bannedAt sql.NullTime

	err := row.Scan(
		&id, &u.Username, &u.Email, &u.PasswordHash,
		&u.Role, &u.Points, &u.CreatedAt, &u.UpdatedAt,
		&u.Banned, &bannedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("user not found")
	}
	if err != nil {
		return nil, err
	}

	u.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}

	if bannedAt.Valid {
		u.BannedAt = &bannedAt.Time
	}

	return &u, nil
}
