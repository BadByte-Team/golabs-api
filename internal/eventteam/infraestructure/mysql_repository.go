// Package infrastructure implementa el repositorio de equipos por evento usando MySQL/MariaDB.
package infrastructure

import (
	"database/sql"
	"errors"
	"time"

	"github.com/google/uuid"

	teamdomain "golabs-api/internal/eventteam/domain"
)

// MySQLEventTeamRepository implementa teamdomain.Repository usando MySQL/MariaDB.
type MySQLEventTeamRepository struct {
	db *sql.DB
}

// NewEventTeamRepository crea una instancia de MySQLEventTeamRepository.
// Retorna la interfaz teamdomain.Repository para desacoplar del tipo concreto.
func NewEventTeamRepository(db *sql.DB) teamdomain.Repository {
	return &MySQLEventTeamRepository{db: db}
}

// SaveTeam inserta un nuevo equipo en la tabla event_teams.
// Establece created_at y updated_at en UTC al momento de la insercion.
func (r *MySQLEventTeamRepository) SaveTeam(team *teamdomain.EventTeam) error {
	query := `
		INSERT INTO event_teams (
			id, event_id, name, join_secret_hash,
			score, created_at, updated_at
		) VALUES (?, ?, ?, ?, ?, ?, ?)
	`

	now := time.Now().UTC()
	team.CreatedAt = now
	team.UpdatedAt = now

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		team.ID.String(),
		team.EventID.String(),
		team.Name,
		team.JoinSecretHash,
		team.Score,
		team.CreatedAt,
		team.UpdatedAt,
	)

	return err
}

// UpdateTeam persiste los cambios en un equipo existente (por ejemplo, rotacion del join secret).
// Actualiza automaticamente el campo updated_at.
func (r *MySQLEventTeamRepository) UpdateTeam(team *teamdomain.EventTeam) error {
	query := `
		UPDATE event_teams
		SET name = ?, join_secret_hash = ?, score = ?, updated_at = ?
		WHERE id = ?
	`

	team.UpdatedAt = time.Now().UTC()

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		team.Name,
		team.JoinSecretHash,
		team.Score,
		team.UpdatedAt,
		team.ID.String(),
	)

	return err
}

// GetTeamByID busca un equipo por su UUID.
func (r *MySQLEventTeamRepository) GetTeamByID(id uuid.UUID) (*teamdomain.EventTeam, error) {
	query := `
		SELECT id, event_id, name, join_secret_hash,
		       score, created_at, updated_at
		FROM event_teams
		WHERE id = ?
		LIMIT 1
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanTeam(smt.QueryRow(id.String()))
}

// GetTeamByName busca un equipo por su nombre dentro de un evento especifico.
// Usado para validar join secret al unirse a un equipo.
func (r *MySQLEventTeamRepository) GetTeamByName(eventID uuid.UUID, name string) (*teamdomain.EventTeam, error) {
	query := `
		SELECT id, event_id, name, join_secret_hash,
		       score, created_at, updated_at
		FROM event_teams
		WHERE event_id = ? AND name = ?
		LIMIT 1
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanTeam(smt.QueryRow(eventID.String(), name))
}

// AddMember inserta un nuevo miembro en la tabla event_team_members.
// Si JoinedAt es zero, se usa el tiempo actual en UTC.
func (r *MySQLEventTeamRepository) AddMember(member *teamdomain.EventTeamMember) error {
	query := `
		INSERT INTO event_team_members (
			event_team_id, user_id, role, joined_at
		) VALUES (?, ?, ?, ?)
	`

	// Usar el JoinedAt del objeto de dominio; default a now si no fue inicializado.
	joinedAt := member.JoinedAt
	if joinedAt.IsZero() {
		joinedAt = time.Now().UTC()
	}

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(
		member.EventTeamID.String(),
		member.UserID.String(),
		string(member.Role),
		joinedAt,
	)

	return err
}

// RemoveMember elimina a un miembro del equipo.
func (r *MySQLEventTeamRepository) RemoveMember(teamID, userID uuid.UUID) error {
	query := `
		DELETE FROM event_team_members
		WHERE event_team_id = ? AND user_id = ?
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return err
	}
	defer smt.Close()

	_, err = smt.Exec(teamID.String(), userID.String())
	return err
}

// ListMembers retorna todos los miembros de un equipo.
func (r *MySQLEventTeamRepository) ListMembers(teamID uuid.UUID) ([]*teamdomain.EventTeamMember, error) {
	query := `
		SELECT event_team_id, user_id, role, joined_at
		FROM event_team_members
		WHERE event_team_id = ?
	`

	rows, err := r.db.Query(query, teamID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var members []*teamdomain.EventTeamMember
	for rows.Next() {
		var m teamdomain.EventTeamMember
		var teamIDStr, userIDStr string
		var role string

		err := rows.Scan(
			&teamIDStr,
			&userIDStr,
			&role,
			&m.JoinedAt,
		)
		if err != nil {
			return nil, err
		}

		m.EventTeamID, err = uuid.Parse(teamIDStr)
		if err != nil {
			return nil, err
		}
		m.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		m.Role = teamdomain.TeamRole(role)

		members = append(members, &m)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return members, nil
}

// CountMembers retorna el numero de miembros activos en el equipo.
// Usado para validar si el equipo tiene espacio antes de agregar un nuevo miembro.
func (r *MySQLEventTeamRepository) CountMembers(teamID uuid.UUID) (int, error) {
	query := `
		SELECT COUNT(*) FROM event_team_members
		WHERE event_team_id = ?
	`

	var count int
	err := r.db.QueryRow(query, teamID.String()).Scan(&count)
	return count, err
}

// IsUserInEvent verifica si el usuario ya pertenece a algun equipo en el evento.
// Usa un JOIN entre event_team_members y event_teams para la verificacion.
func (r *MySQLEventTeamRepository) IsUserInEvent(eventID, userID uuid.UUID) (bool, error) {
	query := `
		SELECT COUNT(*)
		FROM event_team_members etm
		JOIN event_teams et ON et.id = etm.event_team_id
		WHERE et.event_id = ? AND etm.user_id = ?
	`

	var count int
	err := r.db.QueryRow(
		query,
		eventID.String(),
		userID.String(),
	).Scan(&count)

	return count > 0, err
}

// ListTeamsByEvent retorna todos los equipos de un evento ordenados por puntaje descendente.
// El orden por score DESC es el que determina el ranking del leaderboard.
func (r *MySQLEventTeamRepository) ListTeamsByEvent(eventID uuid.UUID) ([]*teamdomain.EventTeam, error) {
	query := `
		SELECT id, event_id, name, join_secret_hash,
		       score, created_at, updated_at
		FROM event_teams
		WHERE event_id = ?
		ORDER BY score DESC
	`
	rows, err := r.db.Query(query, eventID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	teams := make([]*teamdomain.EventTeam, 0)
	for rows.Next() {
		var t teamdomain.EventTeam
		var id, evID string
		err := rows.Scan(&id, &evID, &t.Name, &t.JoinSecretHash, &t.Score, &t.CreatedAt, &t.UpdatedAt)
		if err != nil {
			return nil, err
		}
		t.ID, err = uuid.Parse(id)
		if err != nil {
			return nil, err
		}
		t.EventID, err = uuid.Parse(evID)
		if err != nil {
			return nil, err
		}
		teams = append(teams, &t)
	}
	return teams, rows.Err()
}

// ListMembersWithUsername retorna los miembros de un equipo enriquecidos con el username de cada usuario.
// Realiza un JOIN con la tabla users para resolver el username en una sola consulta.
func (r *MySQLEventTeamRepository) ListMembersWithUsername(teamID uuid.UUID) ([]*teamdomain.MemberWithUsername, error) {
	query := `
		SELECT etm.event_team_id, etm.user_id, u.username, etm.role, etm.joined_at
		FROM event_team_members etm
		JOIN users u ON u.id = etm.user_id
		WHERE etm.event_team_id = ?
		ORDER BY etm.joined_at ASC
	`
	rows, err := r.db.Query(query, teamID.String())
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	members := make([]*teamdomain.MemberWithUsername, 0)
	for rows.Next() {
		var m teamdomain.MemberWithUsername
		var teamIDStr, userIDStr, role string
		err := rows.Scan(&teamIDStr, &userIDStr, &m.Username, &role, &m.JoinedAt)
		if err != nil {
			return nil, err
		}
		m.EventTeamID, err = uuid.Parse(teamIDStr)
		if err != nil {
			return nil, err
		}
		m.UserID, err = uuid.Parse(userIDStr)
		if err != nil {
			return nil, err
		}
		m.Role = teamdomain.TeamRole(role)
		members = append(members, &m)
	}
	return members, rows.Err()
}

// GetTeamByUserAndEvent retorna el equipo al que pertenece un usuario en un evento especifico.
// Usado para verificar la membresia antes de operaciones de equipo sensibles.
func (r *MySQLEventTeamRepository) GetTeamByUserAndEvent(eventID, userID uuid.UUID) (*teamdomain.EventTeam, error) {
	query := `
		SELECT et.id, et.event_id, et.name, et.join_secret_hash,
		       et.score, et.created_at, et.updated_at
		FROM event_teams et
		JOIN event_team_members etm ON etm.event_team_id = et.id
		WHERE et.event_id = ? AND etm.user_id = ?
		LIMIT 1
	`

	smt, err := r.db.Prepare(query)
	if err != nil {
		return nil, err
	}
	defer smt.Close()

	return scanTeam(smt.QueryRow(eventID.String(), userID.String()))
}

// scanTeam mapea un sql.Row a un EventTeam, convirtiendo los UUIDs de string.
func scanTeam(row *sql.Row) (*teamdomain.EventTeam, error) {
	var t teamdomain.EventTeam
	var id, eventID string

	err := row.Scan(
		&id,
		&eventID,
		&t.Name,
		&t.JoinSecretHash,
		&t.Score,
		&t.CreatedAt,
		&t.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, errors.New("team not found")
	}
	if err != nil {
		return nil, err
	}

	t.ID, err = uuid.Parse(id)
	if err != nil {
		return nil, err
	}
	t.EventID, err = uuid.Parse(eventID)
	if err != nil {
		return nil, err
	}

	return &t, nil
}
