// Package infrastructure implementa el repositorio de usuarios usando MySQL/MariaDB.
package infrastructure

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	userdomain "golabs-api/internal/user/domain"
)

// MySQLUserRepository implementa userdomain.UserRepository usando MySQL/MariaDB.
type MySQLUserRepository struct {
	db *sql.DB
}

// NewUserRepository crea una instancia de MySQLUserRepository.
// Retorna la interfaz userdomain.UserRepository para desacoplar del tipo concreto.
func NewUserRepository(db *sql.DB) userdomain.UserRepository {
	return &MySQLUserRepository{db: db}
}

// Create inserta un nuevo usuario en la tabla users.
// Genera un UUID nuevo y establece los timestamps created_at / updated_at en UTC.
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
	defer smt.Close()

	_, err = smt.Exec(
		user.ID.String(), user.Username, user.Email,
		user.PasswordHash, user.Role, user.Points,
		user.CreatedAt, user.UpdatedAt, user.Banned, user.BannedAt,
	)
	return err
}

// GetByEmail busca un usuario por su email. Usado principalmente en el flujo de login.
func (r *MySQLUserRepository) GetByEmail(email string) (*userdomain.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		FROM users WHERE email = ? LIMIT 1
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanUser(smt.QueryRow(email))
}

// GetByID busca un usuario por su UUID (como string).
func (r *MySQLUserRepository) GetByID(id string) (*userdomain.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		FROM users WHERE id = ? LIMIT 1
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanUser(smt.QueryRow(id))
}

// GetByUsername busca un usuario por su username exacto (case-sensitive segun el collation de la BD).
func (r *MySQLUserRepository) GetByUsername(username string) (*userdomain.User, error) {
	query := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		FROM users WHERE username = ? LIMIT 1
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanUser(smt.QueryRow(username))
}

// SearchByUsername busca usuarios cuyo username contenga el termino de busqueda (LIKE %query%).
// Retorna hasta 20 resultados ordenados alfabeticamente.
func (r *MySQLUserRepository) SearchByUsername(query string) ([]*userdomain.User, error) {
	q := `
		SELECT id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		FROM users WHERE username LIKE ? ORDER BY username ASC LIMIT 20
	`
	rows, err := r.db.Query(q, "%"+query+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*userdomain.User, 0)
	for rows.Next() {
		u, err := scanUserRow(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, rows.Err()
}

// Update persiste los cambios en username, email y points de un usuario.
// Actualiza automaticamente el campo updated_at.
func (r *MySQLUserRepository) Update(user *userdomain.User) error {
	query := `UPDATE users SET username = ?, email = ?, points = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()
	user.UpdatedAt = now

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(user.Username, user.Email, user.Points, now, user.ID.String())
	return err
}

// UpdatePassword reemplaza el hash de la contrasena del usuario.
// Solo almacena el hash bcrypt; la contrasena en texto plano no llega aqui.
func (r *MySQLUserRepository) UpdatePassword(userID, passwordHash string) error {
	query := `UPDATE users SET password_hash = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(passwordHash, now, userID)
	return err
}

// UpdateRole cambia el rol del usuario (por ejemplo, de "user" a "admin").
func (r *MySQLUserRepository) UpdateRole(userID, role string) error {
	query := `UPDATE users SET role = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(role, now, userID)
	return err
}

// UpdatePoints reemplaza los puntos totales del usuario.
// Los puntos se definen de forma absoluta (no se suman/restan en la BD).
func (r *MySQLUserRepository) UpdatePoints(userID string, points int) error {
	query := `UPDATE users SET points = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(points, now, userID)
	return err
}

// Ban marca al usuario como baneado y registra la fecha de baneo.
func (r *MySQLUserRepository) Ban(userID string) error {
	query := `UPDATE users SET banned = ?, banned_at = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(true, now, now, userID)
	return err
}

// Unban reactiva a un usuario baneado, limpiando el campo banned_at.
func (r *MySQLUserRepository) Unban(userID string) error {
	query := `UPDATE users SET banned = ?, banned_at = ?, updated_at = ? WHERE id = ?`
	now := time.Now().UTC()

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(false, nil, now, userID)
	return err
}

// List retorna una pagina de usuarios junto con el total disponible.
// Los resultados se ordenan por created_at descendente (mas recientes primero).
func (r *MySQLUserRepository) List(offset, size int) ([]*userdomain.User, int, error) {
	var total int
	if err := r.db.QueryRow(`SELECT COUNT(*) FROM users`).Scan(&total); err != nil {
		return nil, 0, err
	}

	rows, err := r.db.Query(
		`SELECT id, username, email, password_hash, role, points, created_at, updated_at, banned, banned_at
		 FROM users ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		size, offset,
	)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	users := make([]*userdomain.User, 0)
	for rows.Next() {
		u, err := scanUserRow(rows)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, u)
	}
	return users, total, rows.Err()
}

// scanUser mapea un sql.Row a un User. Usado en consultas de fila unica.
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

	return populateUser(&u, id, bannedAt)
}

// scanUserRow mapea un sql.Rows a un User. Usado en consultas de lista.
func scanUserRow(rows *sql.Rows) (*userdomain.User, error) {
	var u userdomain.User
	var id string
	var bannedAt sql.NullTime

	err := rows.Scan(
		&id, &u.Username, &u.Email, &u.PasswordHash,
		&u.Role, &u.Points, &u.CreatedAt, &u.UpdatedAt,
		&u.Banned, &bannedAt,
	)
	if err != nil {
		return nil, err
	}

	return populateUser(&u, id, bannedAt)
}

// populateUser convierte el UUID string y el NullTime de bannedAt al tipo de dominio.
func populateUser(u *userdomain.User, id string, bannedAt sql.NullTime) (*userdomain.User, error) {
	var err error
	u.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	if bannedAt.Valid {
		u.BannedAt = &bannedAt.Time
	}
	return u, nil
}
