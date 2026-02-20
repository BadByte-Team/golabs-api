package http

import (
	"database/sql"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"

	user_uc "golabs-api/internal/application/usecases/user"
	"golabs-api/internal/health"
	repos "golabs-api/internal/infrastructure/db/repositories"
	"golabs-api/internal/infrastructure/security"
	handlers "golabs-api/internal/interfaces/http/handlers"
	authmw "golabs-api/internal/interfaces/http/middleware"
)

func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.Recoverer)

	r.Get("/health", health.NewHandler(db).ServeHTTP)
	userRepo := repos.NewUserRepository(db)
	jwtSvc, err := security.NewJWTService()
	if err != nil {
		panic("JWT no configurado: " + err.Error())
	}

	createUserUC := user_uc.NewCreateUserUseCase(userRepo)
	getUserByIDUC := user_uc.NewGetUserByIDUseCase(userRepo)
	updateUserUC := user_uc.NewUpdateUserUseCase(userRepo)
	changePasswordUC := user_uc.NewChangePasswordUseCase(userRepo)

	updateRoleUC := user_uc.NewUpdateUserRoleUseCase(userRepo)
	updatePointsUC := user_uc.NewUpdateUserPointsUseCase(userRepo)
	banUserUC := user_uc.NewBanUserUseCase(userRepo)
	unbanUserUC := user_uc.NewUnbanUserUseCase(userRepo)

	loginUC := user_uc.NewLoginUseCase(userRepo, jwtSvc)
	authHandler := handlers.NewAuthHandler(loginUC, createUserUC)
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

	r.Route("/auth", func(r chi.Router) {
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
	})

	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuth(jwtSvc))
		r.Use(authmw.LoadUser(userRepo))
		r.Use(authmw.RequireNotBanned)

		r.Route("/users", func(r chi.Router) {

			r.With(authmw.RequireRole("admin")).
				Post("/", userHandler.Create)

			r.Get("/{id}", userHandler.GetByID)

			r.With(authmw.RequireSelfOrAdmin).
				Post("/{id}/update", userHandler.Update)

			r.With(authmw.RequireSelfOrAdmin).
				Post("/{id}/password", userHandler.ChangePassword)
		})

		r.Route("/admin/users", func(r chi.Router) {
			r.Use(authmw.RequireRole("admin"))

			r.Post("/{id}/role", userHandler.UpdateRole)
			r.Post("/{id}/points", userHandler.UpdatePoints)
			r.Post("/{id}/ban", userHandler.Ban)
			r.Post("/{id}/unban", userHandler.Unban)
		})
	})

	return r
}
