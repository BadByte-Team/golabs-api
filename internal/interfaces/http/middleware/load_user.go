package middleware

import (
	"net/http"

	repos "golabs-api/internal/domain/repositories"
)

func LoadUser(userRepo repos.UserRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			ctxUser, ok := GetUser(r.Context())
			if !ok {
				http.Error(w, "unauthorized", http.StatusUnauthorized)
				return
			}

			user, err := userRepo.GetByID(ctxUser.UserID)
			if err != nil {
				http.Error(w, "usuario no encontrado", http.StatusUnauthorized)
				return
			}

			// ⬅️ Fuente de verdad
			ctx := WithUser(r.Context(), UserContext{
				UserID: user.ID.String(),
				Role:   user.Role,
				Banned: user.Banned,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
