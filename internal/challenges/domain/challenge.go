// Package domain define los tipos de dominio del modulo de challenges CTF:
// categorias, niveles de dificultad, el modelo Challenge, Flag y Solve,
// junto con los constructores que aplican las reglas de negocio basicas.
package domain

import (
	"errors"
	"time"

	"github.com/google/uuid"
)

// ChallengeCategory clasifica el tipo de reto segun la disciplina tecnica.
type ChallengeCategory string

const (
	CategoryPwn       ChallengeCategory = "pwn"
	CategoryWeb       ChallengeCategory = "web"
	CategoryCrypto    ChallengeCategory = "crypto"
	CategoryForensics ChallengeCategory = "forensics"
	CategoryReverse   ChallengeCategory = "reverse"
	CategoryOSINT     ChallengeCategory = "osint"
	CategoryMisc      ChallengeCategory = "misc"
)

// ChallengeDifficulty indica el nivel de dificultad estimado del reto.
type ChallengeDifficulty string

const (
	DifficultyEasy   ChallengeDifficulty = "easy"
	DifficultyMedium ChallengeDifficulty = "medium"
	DifficultyHard   ChallengeDifficulty = "hard"
)

// Challenge es el reto CTF publicado dentro de un evento.
// El campo Visible controla si los participantes pueden verlo;
// se cambia a true cuando el admin llama a Publish().
type Challenge struct {
	ID          uuid.UUID
	EventID     uuid.UUID
	Title       string
	Description string
	Category    ChallengeCategory
	Points      int
	Difficulty  ChallengeDifficulty
	FileURL     string // URL para descargar archivo asociado al reto (opcional)
	Visible     bool   // false hasta que el admin publique el challenge explicitamente
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// NewChallenge crea un Challenge en estado oculto (Visible = false).
//
// Reglas:
//   - eventID no puede ser uuid.Nil.
//   - title es obligatorio.
//   - points no puede ser negativo.
//   - category es obligatoria.
//   - si difficulty esta vacia se asigna DifficultyMedium por defecto.
//
// Retorna error si alguna regla se viola.
func NewChallenge(
	eventID uuid.UUID,
	title, description string,
	category ChallengeCategory,
	points int,
	difficulty ChallengeDifficulty,
	fileURL string,
) (*Challenge, error) {
	if eventID == uuid.Nil {
		return nil, errors.New("eventID es requerido")
	}
	if title == "" {
		return nil, errors.New("el titulo del challenge es requerido")
	}
	if points < 0 {
		return nil, errors.New("los puntos no pueden ser negativos")
	}
	if category == "" {
		return nil, errors.New("la categoria es requerida")
	}
	if difficulty == "" {
		// Dificultad por defecto si el admin no especifica.
		difficulty = DifficultyMedium
	}

	now := time.Now()
	return &Challenge{
		ID:          uuid.New(),
		EventID:     eventID,
		Title:       title,
		Description: description,
		Category:    category,
		Points:      points,
		Difficulty:  difficulty,
		FileURL:     fileURL,
		Visible:     false,
		CreatedAt:   now,
		UpdatedAt:   now,
	}, nil
}

// Update modifica los campos editables del challenge.
// No permite cambiar EventID ni Visible; esos se manejan con metodos dedicados.
//
// Retorna error si title esta vacio o points es negativo.
func (c *Challenge) Update(title, description string, category ChallengeCategory, points int, difficulty ChallengeDifficulty, fileURL string) error {
	if title == "" {
		return errors.New("el titulo es requerido")
	}
	if points < 0 {
		return errors.New("los puntos no pueden ser negativos")
	}
	c.Title = title
	c.Description = description
	c.Category = category
	c.Points = points
	c.Difficulty = difficulty
	c.FileURL = fileURL
	c.UpdatedAt = time.Now()
	return nil
}

// Publish hace el challenge visible para los participantes del evento.
func (c *Challenge) Publish() {
	c.Visible = true
	c.UpdatedAt = time.Now()
}

// Unpublish oculta el challenge a los participantes (puede usarse para retirar un reto con error).
func (c *Challenge) Unpublish() {
	c.Visible = false
	c.UpdatedAt = time.Now()
}

// Flag almacena el hash SHA-256 de la flag real de un challenge.
// El valor en texto plano NUNCA se persiste en base de datos.
type Flag struct {
	ID          uuid.UUID
	ChallengeID uuid.UUID
	Hash        string // hash SHA-256 en hexadecimal de la flag real
	CreatedAt   time.Time
}

// Solve es un registro inmutable que se crea cuando un equipo envia la flag correcta.
// Funciona como audit log de la competencia.
type Solve struct {
	ID          uuid.UUID
	ChallengeID uuid.UUID
	EventTeamID uuid.UUID
	UserID      uuid.UUID // usuario especifico del equipo que envio la flag
	SolvedAt    time.Time
}

// NewSolve crea un registro de solve con timestamps y UUID generados automaticamente.
//
// Entradas: challengeID, teamID y userID del solve.
// Salida:    puntero a Solve listo para persistir.
func NewSolve(challengeID, teamID, userID uuid.UUID) *Solve {
	return &Solve{
		ID:          uuid.New(),
		ChallengeID: challengeID,
		EventTeamID: teamID,
		UserID:      userID,
		SolvedAt:    time.Now(),
	}
}
