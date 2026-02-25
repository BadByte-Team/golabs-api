// Package domain define los tipos de dominio del modulo de equipos de evento,
// incluyendo EventTeam, EventTeamMember, roles de equipo y la estructura
// LeaderboardEntry usada para el ranking de puntuaciones.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// TeamRole define el rol que puede tener un miembro dentro de un equipo.
type TeamRole string

const (
	TeamOwner  TeamRole = "owner"  // creador del equipo; puede rotar el secreto de union
	TeamMember TeamRole = "member" // miembro regular
)

// EventTeam representa un equipo inscrito en un evento CTF.
// El secreto de union se almacena hasheado (SHA-256) para no exponerlo en BD.
type EventTeam struct {
	ID             uuid.UUID
	EventID        uuid.UUID
	Name           string
	JoinSecretHash string // hash SHA-256 del secreto de union; nunca el valor plano
	Score          int    // puntos acumulados por el equipo al resolver challenges
	CreatedAt      time.Time
	UpdatedAt      time.Time
}

// EventTeamMember registra la pertenencia de un usuario a un equipo de evento.
type EventTeamMember struct {
	EventTeamID uuid.UUID
	UserID      uuid.UUID
	Role        TeamRole
	JoinedAt    time.Time
}

// LeaderboardEntry es una proyeccion de lectura usada para el ranking del evento.
// Se construye en el use case y nunca se persiste directamente.
type LeaderboardEntry struct {
	Rank        int
	TeamID      string
	TeamName    string
	Score       int
	MemberCount int
}

// NewEventTeam crea un EventTeam aplicando las reglas minimas de validacion.
//
// Reglas:
//   - eventID no puede ser uuid.Nil.
//   - name es obligatorio.
//   - joinSecretHash es obligatorio (debe ser el hash SHA-256 del secreto generado).
//
// Retorna error si alguna regla se viola.
func NewEventTeam(
	eventID uuid.UUID,
	name string,
	joinSecretHash string,
) (*EventTeam, error) {

	if eventID == uuid.Nil {
		return nil, errors.New("eventID es requerido")
	}

	if name == "" {
		return nil, errors.New("el nombre del equipo es requerido")
	}

	if joinSecretHash == "" {
		return nil, errors.New("el hash del secreto de union es requerido")
	}

	now := time.Now()

	return &EventTeam{
		ID:             uuid.New(),
		EventID:        eventID,
		Name:           name,
		JoinSecretHash: joinSecretHash,
		Score:          0,
		CreatedAt:      now,
		UpdatedAt:      now,
	}, nil
}
