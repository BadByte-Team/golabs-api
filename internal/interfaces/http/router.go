package http

import (
	"database/sql"
	"net/http"

	user_usecases "golabs-api/internal/application/usecases/user"
	"golabs-api/internal/health"
	repos "golabs-api/internal/infrastructure/db/repositories"
	"golabs-api/internal/interfaces/http/handlers"
)

func NewRouter(db *sql.DB) http.Handler {
	mux := http.NewServeMux()

	healthHandler := health.NewHandler(db)
	mux.Handle("/health", healthHandler)

	userRepo := repos.NewUserRepository(db)

	createUserUseCase := user_usecases.NewCreateUserUseCase(userRepo)
	getUserByIDUseCase := user_usecases.NewGetUserByIDUseCase(userRepo)
	updateUserUseCase := user_usecases.NewUpdateUserUseCase(userRepo)

	userHandler := handlers.NewUserHandler(
		createUserUseCase,
		getUserByIDUseCase,
		updateUserUseCase,
	)

	mux.HandleFunc("/users", userHandler.Create)
	mux.HandleFunc("/users/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			userHandler.GetByID(w, r)
		case http.MethodPut:
			userHandler.Update(w, r)
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
		}
	})
	return mux
}
