// Package validate provee una instancia compartida del validador de structs
// y helpers para decodificar y validar el body JSON de peticiones HTTP.
//
// Usa github.com/go-playground/validator/v10 para la validacion basada en struct tags.
// Los mensajes de error de validacion se construyen en ingles porque son parte
// de la respuesta de API (consumida por clientes de distintos idiomas).
package validate

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-playground/validator/v10"
)

// v es la instancia compartida del validador. Se crea una sola vez para
// aprovechar el cache interno de reflection que construye el validador.
var v = validator.New()

// DecodeAndValidate decodifica el body JSON de r en dst y ejecuta la validacion
// de struct tags (validate:"required", "email", "min", etc.).
//
// Entrada:  r, peticion HTTP; dst, puntero al struct destino.
// Salida:   nil si todo es valido, error descriptivo si la decodificacion o validacion falla.
//
// Los errores de validacion se concatenan con "; " para retornar todos los campos
// invalidos en un solo mensaje, facilitando la depuracion del cliente.
func DecodeAndValidate(r *http.Request, dst any) error {
	if err := json.NewDecoder(r.Body).Decode(dst); err != nil {
		return fmt.Errorf("invalid JSON: %w", err)
	}
	if err := v.Struct(dst); err != nil {
		var ve validator.ValidationErrors
		if errors.As(err, &ve) {
			msgs := make([]string, 0, len(ve))
			for _, fe := range ve {
				msgs = append(msgs, fieldError(fe))
			}
			return errors.New(strings.Join(msgs, "; "))
		}
		return err
	}
	return nil
}

// DecodeOnly decodifica el body JSON de r en dst sin ejecutar validacion de struct tags.
// Usar cuando la validacion se hace manualmente en el use case, o cuando el struct
// no tiene tags de validacion.
//
// Entrada:  r, peticion HTTP; dst, puntero al struct destino.
// Salida:   error de decodificacion JSON, o nil si fue exitoso.
func DecodeOnly(r *http.Request, dst any) error {
	return json.NewDecoder(r.Body).Decode(dst)
}

// fieldError construye un mensaje de error legible por el cliente
// a partir de un FieldError de go-playground/validator.
func fieldError(fe validator.FieldError) string {
	switch fe.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fe.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fe.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fe.Field(), fe.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fe.Field(), fe.Param())
	case "oneof":
		return fmt.Sprintf("%s must be one of [%s]", fe.Field(), fe.Param())
	case "gt":
		return fmt.Sprintf("%s must be greater than %s", fe.Field(), fe.Param())
	default:
		return fmt.Sprintf("%s failed validation (%s)", fe.Field(), fe.Tag())
	}
}
