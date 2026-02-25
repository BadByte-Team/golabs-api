// Package apperrors defines sentinel domain errors and HTTP mapping helpers.
// Use these in use cases to signal the kind of failure, and in handlers to
// map them to the correct HTTP status code.
package apperrors

import (
	"encoding/json"
	"errors"
	"net/http"
)

// Sentinel errors — use errors.Is() to check.
var (
	ErrNotFound     = errors.New("not found")
	ErrConflict     = errors.New("conflict")
	ErrUnauthorized = errors.New("unauthorized")
	ErrForbidden    = errors.New("forbidden")
	ErrBadRequest   = errors.New("bad request")
)

// httpStatus maps a sentinel error to its HTTP status code.
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

// ErrorResponse is the standard JSON error body.
type ErrorResponse struct {
	Error string `json:"error"`
}

// RespondError writes the correct HTTP status + JSON body for a domain error.
func RespondError(w http.ResponseWriter, err error) {
	status := httpStatus(err)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(ErrorResponse{Error: err.Error()})
}

// RespondJSON writes a successful JSON response with the given status.
func RespondJSON(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}
