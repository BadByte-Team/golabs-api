// Package security agrupa las utilidades criptograficas del servidor:
// generacion y validacion de JWT de acceso, generacion de refresh tokens
// y calculo de hashes para comparacion de flags y secretos de equipo.
package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// BcryptCost es el factor de trabajo usado al hashear contrasenas con bcrypt.
// Incrementar este valor incrementa el costo computacional del hash,
// dificultando ataques de fuerza bruta a medida que el hardware mejora.
const BcryptCost = 12

// joinSecretCharset define el alfabeto de los secretos de union de equipos.
// Solo caracteres alfanumericos para evitar ambiguedad y problemas de URL encoding.
const joinSecretCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Hash retorna el digest SHA-256 en hexadecimal de la cadena s.
//
// Usos: comparacion de flags enviadas por los equipos y verificacion de
// secretos de union de equipos. Nunca se almacena el valor original.
//
// Entrada:  cadena plana (flag, secreto, etc.)
// Salida:   representacion hexadecimal del hash SHA-256 (64 caracteres).
func Hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// HashJoinSecret es un alias de Hash mantenido por compatibilidad con versiones anteriores.
// Preferir Hash directamente en codigo nuevo.
func HashJoinSecret(secret string) string { return Hash(secret) }

// GenerateJoinSecret crea un secreto aleatorio criptograficamente seguro de la longitud indicada.
//
// Entrada:  length, numero de caracteres del secreto generado.
// Salida:   cadena aleatoria del alfabeto [a-zA-Z0-9] o error si falla el CSPRNG del sistema.
//
// Se usa crypto/rand en lugar de math/rand para garantizar imprevisibilidad.
func GenerateJoinSecret(length int) (string, error) {
	secret := make([]byte, length)

	for i := range secret {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(joinSecretCharset))))
		if err != nil {
			return "", err
		}
		secret[i] = joinSecretCharset[n.Int64()]
	}

	return string(secret), nil
}
