// Package apperrors define los errores centinela del dominio y los helpers HTTP
// para traducirlos a respuestas JSON con el codigo de estado correcto.
//
// Uso en use cases: retornar un error centinela (p.ej. ErrNotFound) o un error
// que lo envuelva (fmt.Errorf("%w: ...", apperrors.ErrNotFound)).
// Uso en handlers: llamar RespondError(w, err) y dejar que el mapeo sea automatico.
package apperrors

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Errores centinela del dominio. Usar errors.Is() para comparar,
// ya que los use cases pueden envolver estos errores con contexto adicional.
var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
)

// httpStatus mapea un error centinela al codigo HTTP correspondiente.
// Si el error no coincide con ningun centinela conocido se retorna 400 Bad Request.
func httpStatus(err error) int {
	switch {
	case errors.Is(err, ErrNotFound):
		return http.StatusNotFound
	case errors.Is(err, ErrConflict):
		return http.StatusConflict
	case errors.Is(err, ErrUnauthorized):
		return http.StatusUnauthorized
	case errors.Is(err, ErrForbidden):
		return http.StatusForbidden
	default:
		return http.StatusBadRequest
	}
}

// ErrorResponse es el cuerpo JSON estandar para respuestas de error.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondError escribe el codigo HTTP correcto y el cuerpo JSON de error para un error de dominio.
//
// Entrada:  w, el ResponseWriter de la peticion; err, el error del dominio o de aplicacion.
// Salida:   respuesta HTTP con Content-Type: application/json y el mensaje de error.
func RespondError(w http.ResponseWriter, err error) {
	status := httpStatus(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

// RespondJSON escribe una respuesta JSON exitosa con el codigo de estado indicado.
//
// Entrada:  w, el ResponseWriter; status, codigo HTTP (p.ej. 200, 201); v, cualquier valor serializable.
// Salida:   respuesta HTTP con Content-Type: application/json y el cuerpo codificado en JSON.
func RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
