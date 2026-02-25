// Package domain define los tipos de dominio del modulo de challenges CTF:
// categorias, niveles de dificultad, el modelo Challenge, Flag y Solve,
// junto con los constructores que aplican las reglas de negocio basicas.
package domain

import "github.com/google/uuid"

// Repository define los metodos de persistencia del modulo de challenges.
// Agrupa operaciones sobre Challenge, Flag y Solve en una sola interfaz
// porque los tres estan estrechamente acoplados en el flujo de submit.
type Repository interface {

	// SaveChallenge persiste un nuevo challenge.
	SaveChallenge(c *Challenge) error

	// UpdateChallenge persiste los cambios en un challenge existente.
	UpdateChallenge(c *Challenge) error

	// GetChallengeByID busca un challenge por su UUID. Retorna error si no existe.
	GetChallengeByID(id uuid.UUID) (*Challenge, error)

	// ListChallengesByEvent retorna challenges de un evento con filtros opcionales.
	//   visibleOnly: true para participantes (solo publicos), false para admins (todos).
	//   category:    filtro por categoria; cadena vacia = sin filtro.
	//   difficulty:  filtro por dificultad; cadena vacia = sin filtro.
	ListChallengesByEvent(eventID uuid.UUID, visibleOnly bool, category, difficulty string) ([]*Challenge, error)

	// UpsertFlag inserta o reemplaza la flag de un challenge (1 flag por challenge).
	// Solo el hash SHA-256 se almacena; nunca el valor en texto plano.
	UpsertFlag(f *Flag) error

	// GetFlagByChallengeID retorna la flag asociada a un challenge.
	GetFlagByChallengeID(challengeID uuid.UUID) (*Flag, error)

	// SaveSolve registra que un equipo resolvio un challenge.
	SaveSolve(s *Solve) error

	// HasTeamSolved retorna true si el equipo ya resolvio el challenge.
	// Previene doble submission y doble acreditacion de puntos.
	HasTeamSolved(challengeID, teamID uuid.UUID) (bool, error)

	// ListSolvesByChallenge retorna todos los solves de un challenge.
	ListSolvesByChallenge(challengeID uuid.UUID) ([]*Solve, error)

	// ListSolvesByTeam retorna todos los solves de un equipo.
	ListSolvesByTeam(teamID uuid.UUID) ([]*Solve, error)

	// GetSolveCount retorna cuantos equipos han resuelto un challenge.
	// Util para mostrar estadisticas de dificultad dinamica.
	GetSolveCount(challengeID uuid.UUID) (int, error)

	// GetFirstBlood retorna el primer solve registrado para un challenge (first blood).
	// Retorna nil si nadie lo ha resuelto todavia.
	GetFirstBlood(challengeID uuid.UUID) (*Solve, error)
}
