// Package ratelimit implementa un limitador de tasa en memoria basado en una ventana deslizante.
// Se usa internamente por los middlewares LoginRateLimit y UserRateLimit.
//
// Nota: este limitador es por instancia de proceso. En despliegues con multiples instancias
// (horizontal scaling) es necesario usar un limitador distribuido (Redis, etc.).
package ratelimit

import (
	"sync"
	"time"
)

// rateLimiter implementa un contador de ventana deslizante protegido con mutex.
// Cada key (IP, userID, etc.) tiene su propio historial de timestamps.
type rateLimiter struct {
	mu       sync.Mutex
	visits   map[string][]time.Time // historial de timestamps por key
	limit    int                    // numero maximo de visitas permitidas en el intervalo
	interval time.Duration          // duracion de la ventana de tiempo
}

// newRateLimiter crea un nuevo rate limiter con el limite y el intervalo indicados.
//
// Entrada:  limit, numero maximo de requests; interval, ventana de tiempo.
// Salida:   puntero al rateLimiter inicializado.
func newRateLimiter(limit int, interval time.Duration) *rateLimiter {
	return &rateLimiter{
		visits:   make(map[string][]time.Time),
		limit:    limit,
		interval: interval,
	}
}

// allow verifica si la key puede realizar una peticion adicional dentro de la ventana actual.
//
// Limpia los timestamps fuera de la ventana antes de evaluar, lo que implementa
// el comportamiento de ventana deslizante.
//
// Entrada:  key, identificador unico del cliente (IP, userID, etc.)
// Salida:   true si el request esta permitido; false si el limite fue alcanzado.
func (rl *rateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	window := now.Add(-rl.interval)
	times := rl.visits[key]
	valid := times[:0] // reutiliza el slice subyacente para evitar allocations

	// Filtrar timestamps anteriores al inicio de la ventana actual.
	for _, t := range times {
		if t.After(window) {
			valid = append(valid, t)
		}
	}

	if len(valid) >= rl.limit {
		rl.visits[key] = valid
		return false // limite alcanzado
	}

	valid = append(valid, now)
	rl.visits[key] = valid
	return true
}
