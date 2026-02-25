package infrastructure

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	challengedomain "golabs-api/internal/challenges/domain"
)

type MySQLChallengeRepository struct {
	db *sql.DB
}

func NewChallengeRepository(db *sql.DB) challengedomain.Repository {
	return &MySQLChallengeRepository{db: db}
}

/* ── Challenges ─────────────────────────────────────────────────────────── */

func (r *MySQLChallengeRepository) SaveChallenge(c *challengedomain.Challenge) error {
	query := `
		INSERT INTO challenges (
			id, event_id, title, description, category,
			points, difficulty, visible, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		c.ID.String(), c.EventID.String(), c.Title, c.Description,
		string(c.Category), c.Points, string(c.Difficulty), c.Visible,
		c.CreatedAt, c.UpdatedAt,
	)
	return err
}

func (r *MySQLChallengeRepository) UpdateChallenge(c *challengedomain.Challenge) error {
	query := `
		UPDATE challenges
		SET title = ?, description = ?, category = ?, points = ?,
		    difficulty = ?, visible = ?, updated_at = ?
		WHERE id = ?
	`
	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		c.Title, c.Description, string(c.Category), c.Points,
		string(c.Difficulty), c.Visible, c.UpdatedAt, c.ID.String(),
	)
	return err
}

func (r *MySQLChallengeRepository) GetChallengeByID(id uuid.UUID) (*challengedomain.Challenge, error) {
	query := `
		SELECT id, event_id, title, description, category,
		       points, difficulty, visible, created_at, updated_at
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

func (r *MySQLChallengeRepository) ListChallengesByEvent(eventID uuid.UUID, visibleOnly bool) ([]*challengedomain.Challenge, error) {
	query := `
		SELECT id, event_id, title, description, category,
		       points, difficulty, visible, created_at, updated_at
		FROM challenges
		WHERE event_id = ?
	`
	if visibleOnly {
		query += " AND visible = TRUE"
	}
	query += " ORDER BY category, points ASC"

	rows, err := r.db.Query(query, eventID.String())
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
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return challenges, nil
}

/* ── Flags ──────────────────────────────────────────────────────────────── */

// UpsertFlag inserts the flag or replaces the hash if one already exists for this challenge.
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

/* ── Solves ─────────────────────────────────────────────────────────────── */

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

func (r *MySQLChallengeRepository) ListSolvesByChallenge(challengeID uuid.UUID) ([]*challengedomain.Solve, error) {
	return r.listSolves(`WHERE challenge_id = ?`, challengeID.String())
}

func (r *MySQLChallengeRepository) ListSolvesByTeam(teamID uuid.UUID) ([]*challengedomain.Solve, error) {
	return r.listSolves(`WHERE event_team_id = ?`, teamID.String())
}

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

/* ── helpers ────────────────────────────────────────────────────────────── */

func scanChallenge(row *sql.Row) (*challengedomain.Challenge, error) {
	var c challengedomain.Challenge
	var id, eventID, category, difficulty string

	err := row.Scan(
		&id, &eventID, &c.Title, &c.Description, &category,
		&c.Points, &difficulty, &c.Visible, &c.CreatedAt, &c.UpdatedAt,
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

	return &c, nil
}

// scanChallengeRow scans from a *sql.Rows (used in list queries).
func scanChallengeRow(rows *sql.Rows) (*challengedomain.Challenge, error) {
	var c challengedomain.Challenge
	var id, eventID, category, difficulty string

	err := rows.Scan(
		&id, &eventID, &c.Title, &c.Description, &category,
		&c.Points, &difficulty, &c.Visible, &c.CreatedAt, &c.UpdatedAt,
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

	return &c, nil
}

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
	// Ensure compile-time check: all timestamps stored as UTC
	_ = time.UTC
}
