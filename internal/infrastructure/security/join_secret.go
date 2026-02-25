package security

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"math/big"
)

// BcryptCost is the work factor used when hashing passwords.
// Adjust upward over time as hardware gets faster.
const BcryptCost = 12

const joinSecretCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

// Hash returns the SHA-256 hex digest of s.
// Used for flag and join-secret hashing.
func Hash(s string) string {
	sum := sha256.Sum256([]byte(s))
	return hex.EncodeToString(sum[:])
}

// HashJoinSecret is an alias of Hash kept for backward compatibility.
func HashJoinSecret(secret string) string { return Hash(secret) }

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
