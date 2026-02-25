// Package db provee helpers para leer variables de entorno del pool de conexiones.
// Estas funciones retornan el valor por defecto si la variable no esta definida
// o no puede parsearse, evitando panics en arranque.
package db

import (
	"os"
	"strconv"
	"time"
)

// getEnvInt lee la variable de entorno key como entero.
// Si no existe o falla el parseo retorna defaultVal.
func getEnvInt(key string, defaultVal int) int {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return defaultVal
	}
	return parsed
}

// getEnvDuration lee la variable de entorno key como time.Duration (formato Go, ej. "5m").
// Si no existe o falla el parseo retorna defaultVal.
func getEnvDuration(key string, defaultVal time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return defaultVal
	}

	parsed, err := time.ParseDuration(val)
	if err != nil {
		return defaultVal
	}
	return parsed
}
