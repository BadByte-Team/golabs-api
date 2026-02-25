// Package authmw provee las primitivas de contexto HTTP para transportar
// la identidad del usuario autenticado a traves de la cadena de middleware.
package authmw

import (
	"net/http"

	userdomain "golabs-api/internal/user/domain"
)

// LoadUser es un middleware que enriquece el UserContext con datos frescos de la base de datos.
//
// Se usa despues de JWTAuth para rellenar el campo Banned (que no esta en el JWT) y asegurar
// que el usuario sigue existiendo en BD. Si el usuario fue eliminado o baneado desde que emitio
// el token, este middleware rechaza la peticion antes de llegar al handler.
//
// Flujo:
//  1. Extrae el UserContext del contexto (puesto por JWTAuth).
//  2. Busca el usuario en BD por su ID.
//  3. Actualiza el UserContext con el rol actual y el estado de baneo.
//  4. Reemplaza el contexto y pasa al siguiente handler.
//
// Rechaza la peticion con HTTP 401 si no hay UserContext o si el usuario no existe en BD.
func LoadUser(userRepo userdomain.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctxUser, ok := GetUser(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(ctxUser.UserID)
			if err != nil {
				// El usuario puede haber sido eliminado tras emitir el token.
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			// Actualizar con datos frescos de BD, incluyendo estado de baneo.
			ctx := WithUser(r.Context(), UserContext{
				UserID: user.ID.String(),
				Role:   user.Role,
				Banned: user.Banned,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
