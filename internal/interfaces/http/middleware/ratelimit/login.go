// Package ratelimit implementa un limitador de tasa en memoria basado en una ventana deslizante.
package ratelimit

import (
	"net"
	"net/http"
	"strings"
	"time"
)

// loginLimiter permite como maximo 5 intentos de login por IP+email por minuto.
// Previene ataques de fuerza bruta y credential stuffing.
var loginLimiter = newRateLimiter(5, time.Minute)

// LoginRateLimit limita los intentos de autenticacion por combinacion de IP y email.
//
// La clave es "IP:email" (email en minusculas) para distinguir ataques por diccionario
// contra distintas cuentas desde la misma IP.
//
// Retorna HTTP 429 Too Many Requests si se supera el limite.
func LoginRateLimit(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := clientIP(r)
		email := r.URL.Query().Get("email")
		key := ip + ":" + strings.ToLower(email)

		if !loginLimiter.allow(key) {
			http.Error(w, "demasiados intentos, intenta mas tarde", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP extrae la IP real del cliente del campo RemoteAddr.
// Si el formato no es "host:port", retorna RemoteAddr completo como fallback.
func clientIP(r *http.Request) string {
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}
