// Package authmw provee las primitivas de contexto HTTP para transportar
// la identidad del usuario autenticado a traves de la cadena de middleware.
//
// El tipo UserContext se almacena en el contexto de cada peticion autenticada
// con la clave privada userKey para evitar colisiones con otros paquetes.
package authmw

import "context"

// UserContext contiene los datos del usuario autenticado extraidos del JWT.
// Se almacena en el contexto de la peticion tras validar el token.
type UserContext struct {
	UserID string // UUID del usuario en formato string (claim "sub" del JWT)
	Role   string // rol del usuario: "admin" | "user" (claim "role" del JWT)
	Banned bool   // true si el usuario esta baneado; poblado por LoadUser
}

// ctxKey es un tipo privado para la clave de contexto.
// Usar un tipo propio previene colisiones con otras librerias que usen string como clave.
type ctxKey string

// userKey es la clave bajo la que se almacena UserContext en el contexto HTTP.
const userKey ctxKey = "user"

// WithUser retorna un nuevo contexto con el UserContext adjunto.
// Llamado por los middlewares JWT y LoadUser al autenticar una peticion.
//
// Entrada:  ctx, contexto base; user, datos del usuario autenticado.
// Salida:   nuevo contexto que contiene el UserContext.
func WithUser(ctx context.Context, user UserContext) context.Context {
	return context.WithValue(ctx, userKey, user)
}

// GetUser extrae el UserContext del contexto de la peticion.
//
// Salida: el UserContext y true si existe; zero value y false si no hay usuario autenticado.
// Los handlers deben verificar el segundo valor antes de usar el UserContext.
func GetUser(ctx context.Context) (UserContext, bool) {
	user, ok := ctx.Value(userKey).(UserContext)
	return user, ok
}
