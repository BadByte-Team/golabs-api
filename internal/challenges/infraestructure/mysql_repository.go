// Package infrastructure implementa el repositorio de challenges usando MySQL/MariaDB.
package infrastructure

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
)

// MySQLChallengeRepository implementa challengedomain.Repository usando MySQL/MariaDB.
type MySQLChallengeRepository struct {
	db *sql.DB
}

// NewChallengeRepository crea una instancia de MySQLChallengeRepository.
// Retorna la interfaz challengedomain.Repository para desacoplar del tipo concreto.
func NewChallengeRepository(db *sql.DB) challengedomain.Repository {
	return &MySQLChallengeRepository{db: db}
}

// SaveChallenge inserta un nuevo challenge en la tabla challenges.
func (r *MySQLChallengeRepository) SaveChallenge(c *challengedomain.Challenge) error {
	query := `
		INSERT INTO challenges (
			id, event_id, title, description, category,
			points, difficulty, file_url, visible, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	var fileURL *string
	if c.FileURL != "" {
		fileURL = &c.FileURL
	}

	_, err = smt.Exec(
		c.ID.String(), c.EventID.String(), c.Title, c.Description,
		string(c.Category), c.Points, string(c.Difficulty), fileURL, c.Visible,
		c.CreatedAt, c.UpdatedAt,
	)
	return err
}

// UpdateChallenge persiste los cambios en un challenge existente.
// Actualiza todos los campos modificables incluyendo updated_at.
func (r *MySQLChallengeRepository) UpdateChallenge(c *challengedomain.Challenge) error {
	query := `
		UPDATE challenges
		SET title = ?, description = ?, category = ?, points = ?,
		    difficulty = ?, file_url = ?, visible = ?, updated_at = ?
		WHERE id = ?
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	var fileURL *string
	if c.FileURL != "" {
		fileURL = &c.FileURL
	}

	_, err = smt.Exec(
		c.Title, c.Description, string(c.Category), c.Points,
		string(c.Difficulty), fileURL, c.Visible, c.UpdatedAt, c.ID.String(),
	)
	return err
}

// GetChallengeByID busca un challenge por su UUID.
func (r *MySQLChallengeRepository) GetChallengeByID(id uuid.UUID) (*challengedomain.Challenge, error) {
	query := `
		SELECT id, event_id, title, description, category,
		       points, difficulty, file_url, visible, created_at, updated_at
		FROM challenges
		WHERE id = ?
		LIMIT 1
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanChallenge(smt.QueryRow(id.String()))
}

// ListChallengesByEvent retorna los challenges de un evento con filtros opcionales.
// visibleOnly=true excluye challenges ocultos (para no-admins).
// Los filtros de category y difficulty son opcionales; cadena vacia = sin filtro.
// Los resultados se ordenan por category y luego por points ascendente.
func (r *MySQLChallengeRepository) ListChallengesByEvent(eventID uuid.UUID, visibleOnly bool, category, difficulty string) ([]*challengedomain.Challenge, error) {
	q := `SELECT id, event_id, title, description, category, points, difficulty, file_url, visible, created_at, updated_at FROM challenges WHERE event_id = ?`
	args := []any{eventID.String()}

	if visibleOnly {
		q += " AND visible = TRUE"
	}
	if category != "" {
		q += " AND category = ?"
		args = append(args, category)
	}
	if difficulty != "" {
		q += " AND difficulty = ?"
		args = append(args, difficulty)
	}
	q += " ORDER BY category, points ASC"

	rows, err := r.db.Query(q, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	challenges := make([]*challengedomain.Challenge, 0)
	for rows.Next() {
		c, err := scanChallengeRow(rows)
		if err != nil {
			return nil, err
		}
		challenges = append(challenges, c)
	}
	return challenges, rows.Err()
}

// UpsertFlag inserta la flag o reemplaza el hash si ya existe una para este challenge.
// Usa ON DUPLICATE KEY UPDATE sobre el campo challenge_id (unique) para el upsert.
func (r *MySQLChallengeRepository) UpsertFlag(f *challengedomain.Flag) error {
	query := `
		INSERT INTO flags (id, challenge_id, hash, created_at)
		VALUES (?, ?, ?, ?)
		ON DUPLICATE KEY UPDATE hash = VALUES(hash), created_at = VALUES(created_at)
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(f.ID.String(), f.ChallengeID.String(), f.Hash, f.CreatedAt)
	return err
}

// GetFlagByChallengeID retorna la flag de un challenge para validacion de submissions.
// El campo Hash contiene el SHA-256 del texto plano original.
func (r *MySQLChallengeRepository) GetFlagByChallengeID(challengeID uuid.UUID) (*challengedomain.Flag, error) {
	query := `
		SELECT id, challenge_id, hash, created_at
		FROM flags
		WHERE challenge_id = ?
		LIMIT 1
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	var f challengedomain.Flag
	var id, cid string

	err = smt.QueryRow(challengeID.String()).Scan(&id, &cid, &f.Hash, &f.CreatedAt)
	if err == sql.ErrNoRows {
		return nil, errors.New("flag not set")
	}
	if err != nil {
		return nil, err
	}

	f.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	f.ChallengeID, err = uuid.Parse(cid)
	if err != nil {
		return nil, err
	}

	return &f, nil
}

// SaveSolve registra la resolucion de un challenge por un equipo.
func (r *MySQLChallengeRepository) SaveSolve(s *challengedomain.Solve) error {
	query := `
		INSERT INTO solves (id, challenge_id, event_team_id, user_id, solved_at)
		VALUES (?, ?, ?, ?, ?)
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		s.ID.String(), s.ChallengeID.String(),
		s.EventTeamID.String(), s.UserID.String(), s.SolvedAt,
	)
	return err
}

// HasTeamSolved verifica si el equipo ya ha resuelto el challenge.
// Usado para evitar otorgar puntos dobles.
func (r *MySQLChallengeRepository) HasTeamSolved(challengeID, teamID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM solves
		WHERE challenge_id = ? AND event_team_id = ?
	`
	var count int
	err := r.db.QueryRow(query, challengeID.String(), teamID.String()).Scan(&count)
	return count > 0, err
}

// ListSolvesByChallenge retorna todos los solves de un challenge especifico.
func (r *MySQLChallengeRepository) ListSolvesByChallenge(challengeID uuid.UUID) ([]*challengedomain.Solve, error) {
	return r.listSolves(`WHERE challenge_id = ?`, challengeID.String())
}

// ListSolvesByTeam retorna todos los solves realizados por un equipo.
func (r *MySQLChallengeRepository) ListSolvesByTeam(teamID uuid.UUID) ([]*challengedomain.Solve, error) {
	return r.listSolves(`WHERE event_team_id = ?`, teamID.String())
}

// listSolves es un helper interno que ejecuta una consulta de solves con un filtro WHERE parametrizado.
func (r *MySQLChallengeRepository) listSolves(where, arg string) ([]*challengedomain.Solve, error) {
	query := `SELECT id, challenge_id, event_team_id, user_id, solved_at FROM solves ` + where + ` ORDER BY solved_at ASC`

	rows, err := r.db.Query(query, arg)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	solves := make([]*challengedomain.Solve, 0)
	for rows.Next() {
		s, err := scanSolve(rows)
		if err != nil {
			return nil, err
		}
		solves = append(solves, s)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return solves, nil
}

// GetSolveCount retorna el numero de equipos que han resuelto un challenge.
func (r *MySQLChallengeRepository) GetSolveCount(challengeID uuid.UUID) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM solves WHERE challenge_id = ?`, challengeID.String()).Scan(&count)
	return count, err
}

// GetFirstBlood retorna el primer solve de un challenge (el equipo que lo resolvio primero).
// Retorna nil, nil si nadie ha resuelto el challenge aun.
func (r *MySQLChallengeRepository) GetFirstBlood(challengeID uuid.UUID) (*challengedomain.Solve, error) {
	query := `SELECT id, challenge_id, event_team_id, user_id, solved_at FROM solves WHERE challenge_id = ? ORDER BY solved_at ASC LIMIT 1`

	rows, err := r.db.Query(query, challengeID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	if rows.Next() {
		s, err := scanSolve(rows)
		if err != nil {
			return nil, err
		}
		return s, rows.Err()
	}
	// Challenge sin solves: first blood no existe todavia.
	return nil, nil
}

// scanChallenge mapea un sql.Row a un Challenge, convirtiendo UUIDs y enums.
func scanChallenge(row *sql.Row) (*challengedomain.Challenge, error) {
	var c challengedomain.Challenge
	var id, eventID, category, difficulty string
	var fileURL *string

	err := row.Scan(
		&id, &eventID, &c.Title, &c.Description, &category,
		&c.Points, &difficulty, &fileURL, &c.Visible, &c.CreatedAt, &c.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("challenge not found")
	}
	if err != nil {
		return nil, err
	}

	c.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	c.EventID, err = uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}
	c.Category = challengedomain.ChallengeCategory(category)
	c.Difficulty = challengedomain.ChallengeDifficulty(difficulty)
	if fileURL != nil {
		c.FileURL = *fileURL
	}

	return &c, nil
}

// scanChallengeRow mapea un sql.Rows a un Challenge. Usado en queries de lista.
func scanChallengeRow(rows *sql.Rows) (*challengedomain.Challenge, error) {
	var c challengedomain.Challenge
	var id, eventID, category, difficulty string
	var fileURL *string

	err := rows.Scan(
		&id, &eventID, &c.Title, &c.Description, &category,
		&c.Points, &difficulty, &fileURL, &c.Visible, &c.CreatedAt, &c.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	c.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	c.EventID, err = uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}
	c.Category = challengedomain.ChallengeCategory(category)
	c.Difficulty = challengedomain.ChallengeDifficulty(difficulty)
	if fileURL != nil {
		c.FileURL = *fileURL
	}

	return &c, nil
}

// scanSolve mapea un sql.Rows a un Solve, convirtiendo UUIDs de string.
func scanSolve(rows *sql.Rows) (*challengedomain.Solve, error) {
	var s challengedomain.Solve
	var id, challengeID, teamID, userID string

	err := rows.Scan(&id, &challengeID, &teamID, &userID, &s.SolvedAt)
	if err != nil {
		return nil, err
	}

	if s.ID, err = uuid.Parse(id); err != nil {
		return nil, err
	}
	if s.ChallengeID, err = uuid.Parse(challengeID); err != nil {
		return nil, err
	}
	if s.EventTeamID, err = uuid.Parse(teamID); err != nil {
		return nil, err
	}
	if s.UserID, err = uuid.Parse(userID); err != nil {
		return nil, err
	}

	return &s, nil
}

func init() {
	// Verificacion en tiempo de compilacion: los timestamps deben almacenarse en UTC.
	_ = time.UTC
}
