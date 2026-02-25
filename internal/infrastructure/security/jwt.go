package security

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JWTService struct {
	secret   []byte
	issuer   string
	duration time.Duration
}

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
		mins = 15 // short-lived access token; refresh tokens handle long sessions
	}

	return &JWTService{
		secret:   []byte(secret),
		issuer:   issuer,
		duration: time.Duration(mins) * time.Minute,
	}, nil
}

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

func (j *JWTService) Parse(tokenStr string) (*jwt.Token, error) {
	return jwt.Parse(tokenStr, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("método inválido")
		}
		return j.secret, nil
	})
}
