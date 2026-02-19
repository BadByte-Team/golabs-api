package http

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	user_usecases "golabs-api/internal/application/usecases/user"
	"golabs-api/internal/health"
	repos "golabs-api/internal/infrastructure/db/repositories"
	handlers "golabs-api/internal/interfaces/http/handlers"
	authmw "golabs-api/internal/interfaces/http/middleware"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// ======================
	// GLOBAL MIDDLEWARES
	// ======================
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	// TEMPORAL (luego JWT)
	r.Use(authmw.FakeAuth)

	// ======================
	// HEALTH
	// ======================
	r.Get("/health", health.NewHandler(db).ServeHTTP)

	// ======================
	// DEPENDENCIAS
	// ======================
	userRepo := repos.NewUserRepository(db)

	createUserUC := user_usecases.NewCreateUserUseCase(userRepo)
	getUserByIDUC := user_usecases.NewGetUserByIDUseCase(userRepo)
	updateUserUC := user_usecases.NewUpdateUserUseCase(userRepo)
	changePasswordUC := user_usecases.NewChangePasswordUseCase(userRepo)

	updateRoleUC := user_usecases.NewUpdateUserRoleUseCase(userRepo)
	updatePointsUC := user_usecases.NewUpdateUserPointsUseCase(userRepo)
	banUserUC := user_usecases.NewBanUserUseCase(userRepo)
	unbanUserUC := user_usecases.NewUnbanUserUseCase(userRepo)

	userHandler := handlers.NewUserHandler(
		createUserUC,
		getUserByIDUC,
		updateUserUC,
		changePasswordUC,
		updateRoleUC,
		updatePointsUC,
		banUserUC,
		unbanUserUC,
	)

	// ======================
	// USER ROUTES
	// ======================
	r.Route("/users", func(r chi.Router) {

		// Crear usuario → solo admin
		r.With(authmw.RequireRole("admin")).
			Post("/", userHandler.Create)

		// Ver perfil
		r.Get("/{id}", userHandler.GetByID)

		// Update perfil (username/email)
		r.Post("/{id}/update", userHandler.Update)

		// Cambiar contraseña
		r.Post("/{id}/password", userHandler.ChangePassword)
	})

	// ======================
	// ADMIN ROUTES
	// ======================
	r.Route("/admin/users", func(r chi.Router) {
		r.Use(authmw.RequireRole("admin"))

		r.Post("/{id}/role", userHandler.UpdateRole)
		r.Post("/{id}/points", userHandler.UpdatePoints)
		r.Post("/{id}/ban", userHandler.Ban)
		r.Post("/{id}/unban", userHandler.Unban)
	})

	return r
}
