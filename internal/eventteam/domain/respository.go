// Package domain define los tipos de dominio del modulo de equipos de evento,
// incluyendo EventTeam, EventTeamMember, roles de equipo y la estructura
// LeaderboardEntry usada para el ranking de puntuaciones.
package domain

import (
	"time"

	"github.com/google/uuid"
)

// MemberWithUsername enriquece EventTeamMember con el nombre de usuario resuelto en BD.
// Se construye al listar miembros para evitar viajes adicionales a la base de datos en el handler.
type MemberWithUsername struct {
	EventTeamID uuid.UUID
	UserID      uuid.UUID
	Username    string
	Role        TeamRole
	JoinedAt    time.Time
}

// Repository define los metodos de persistencia del agregado EventTeam.
// Las implementaciones concretas viven en el paquete infrastructure.
type Repository interface {
	// SaveTeam persiste un nuevo equipo.
	SaveTeam(team *EventTeam) error

	// UpdateTeam persiste los cambios de un equipo (score, secreto rotado, etc.).
	UpdateTeam(team *EventTeam) error

	// GetTeamByID busca un equipo por su UUID. Retorna error si no existe.
	GetTeamByID(id uuid.UUID) (*EventTeam, error)

	// GetTeamByName busca un equipo por nombre dentro de un evento.
	// Util para validar nombres duplicados al crear equipos.
	GetTeamByName(eventID uuid.UUID, name string) (*EventTeam, error)

	// GetTeamByUserAndEvent retorna el equipo al que pertenece el usuario en el evento dado.
	// Retorna error si el usuario no pertenece a ningun equipo en ese evento.
	GetTeamByUserAndEvent(eventID, userID uuid.UUID) (*EventTeam, error)

	// ListTeamsByEvent retorna todos los equipos de un evento ordenados por score descendente.
	// Usado para el leaderboard y la lista de equipos.
	ListTeamsByEvent(eventID uuid.UUID) ([]*EventTeam, error)

	// AddMember inscribe a un usuario en un equipo con el rol indicado.
	AddMember(member *EventTeamMember) error

	// RemoveMember elimina a un usuario de un equipo.
	RemoveMember(eventTeamID, userID uuid.UUID) error

	// ListMembers retorna los miembros de un equipo sin datos adicionales del usuario.
	ListMembers(eventTeamID uuid.UUID) ([]*EventTeamMember, error)

	// ListMembersWithUsername hace un JOIN con la tabla users para incluir el username
	// de cada miembro, evitando llamadas adicionales desde el handler.
	ListMembersWithUsername(eventTeamID uuid.UUID) ([]*MemberWithUsername, error)

	// CountMembers retorna el numero de miembros de un equipo.
	// Se usa para validar el limite maxTeamSize y en el leaderboard.
	CountMembers(eventTeamID uuid.UUID) (int, error)

	// IsUserInEvent verifica si un usuario ya pertenece a algun equipo en el evento.
	// Previene que un usuario se una a mas de un equipo por evento.
	IsUserInEvent(eventID, userID uuid.UUID) (bool, error)
}
