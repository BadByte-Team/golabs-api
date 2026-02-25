// Package security agrupa las utilidades criptograficas del servidor:
// generacion y validacion de JWT de acceso, generacion de refresh tokens
// y calculo de hashes para comparacion de flags y secretos de equipo.
package security

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// JWTService encapsula la clave secreta, el issuer y la duracion de vida
// de los access tokens firmados con HMAC-SHA256.
type JWTService struct {
	secret   []byte
	issuer   string
	duration time.Duration // duracion del access token (configurable via JWT_EXP_MINUTES)
}

// NewJWTService construye un JWTService leyendo la configuracion desde variables de entorno.
//
// Variables de entorno:
//   - JWT_SECRET (requerido): clave HMAC usada para firmar y verificar tokens.
//   - JWT_ISSUER  (opcional): claim "iss" del token; por defecto "golabs-api".
//   - JWT_EXP_MINUTES (opcional): duracion del access token en minutos; por defecto 15.
//
// Retorna error si JWT_SECRET no esta definido.
func NewJWTService() (*JWTService, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, errors.New("JWT_SECRET no definido")
	}

	issuer := os.Getenv("JWT_ISSUER")
	if issuer == "" {
		issuer = "golabs-api"
	}

	mins, _ := strconv.Atoi(os.Getenv("JWT_EXP_MINUTES"))
	if mins <= 0 {
		// Token de corta duracion; el refresh token gestiona sesiones largas.
		mins = 15
	}

	return &JWTService{
		secret:   []byte(secret),
		issuer:   issuer,
		duration: time.Duration(mins) * time.Minute,
	}, nil
}

// Generate emite un JWT firmado para el usuario indicado.
//
// Entrada:
//   - userID: UUID del usuario como string (claim "sub").
//   - role:   rol del usuario, por ejemplo "admin" o "user" (claim "role").
//
// Salida: token firmado listo para incluir en el header Authorization: Bearer.
func (j *JWTService) Generate(userID, role string) (string, error) {
	now := time.Now()

	claims := jwt.MapClaims{
		"sub":  userID,
		"role": role,
		"iss":  j.issuer,
		"iat":  now.Unix(),
		"exp":  now.Add(j.duration).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(j.secret)
}

// Parse valida la firma del token y retorna el objeto jwt.Token con los claims.
//
// Retorna error si el token esta expirado, la firma es invalida o el algoritmo
// de firma no es HMAC (prevencion de ataques de algoritmo "none").
func (j *JWTService) Parse(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		// Solo se acepta HMAC para evitar el ataque de sustitucion de algoritmo.
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("metodo de firma invalido")
		}
		return j.secret, nil
	})
}
