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

	userHandler := handlers.NewUserHandler(
		createUserUseCase,
		getUserByIDUseCase,
	)

	mux.HandleFunc("/users", userHandler.Create)
	mux.HandleFunc("/users/", userHandler.GetByID)

	return mux
}
