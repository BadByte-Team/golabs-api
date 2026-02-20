package userhttp

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	"golabs-api/internal/infrastructure/security"
	accessmw "golabs-api/internal/interfaces/http/middleware/access"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/middleware/ratelimit"
	userapp "golabs-api/internal/user/application"
	userinfra "golabs-api/internal/user/infrastructure"
)

// RegisterRoutes wires all user and auth routes into the given router.
// It builds the full dependency graph internally so router.go stays clean.
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {
	repo := userinfra.NewUserRepository(db)

	// Use cases
	createUserUC := userapp.NewCreateUserUseCase(repo)
	getUserByIDUC := userapp.NewGetUserByIDUseCase(repo)
	updateUserUC := userapp.NewUpdateUserUseCase(repo)
	changePasswordUC := userapp.NewChangePasswordUseCase(repo)
	updateRoleUC := userapp.NewUpdateUserRoleUseCase(repo)
	updatePointsUC := userapp.NewUpdateUserPointsUseCase(repo)
	banUserUC := userapp.NewBanUserUseCase(repo)
	unbanUserUC := userapp.NewUnbanUserUseCase(repo)
	loginUC := userapp.NewLoginUseCase(repo, jwtSvc)

	// Handlers
	authHandler := NewAuthHandler(loginUC, createUserUC)
	userHandler := NewUserHandler(
		createUserUC, getUserByIDUC, updateUserUC, changePasswordUC,
		updateRoleUC, updatePointsUC, banUserUC, unbanUserUC,
	)

	// Public auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(ratelimit.LoginRateLimit)
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
	})

	// Protected routes
	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuth(jwtSvc))
		r.Use(authmw.LoadUser(repo))
		r.Use(accessmw.RequireNotBanned)
		r.Use(ratelimit.UserRateLimit)

		r.Route("/users", func(r chi.Router) {
			r.With(accessmw.RequireRole("admin")).Post("/", userHandler.Create)
			r.Get("/{id}", userHandler.GetByID)
			r.With(accessmw.RequireSelfOrAdmin).Post("/{id}/update", userHandler.Update)
			r.With(accessmw.RequireSelfOrAdmin).Post("/{id}/password", userHandler.ChangePassword)
		})

		r.Route("/admin/users", func(r chi.Router) {
			r.Use(accessmw.RequireRole("admin"))
			r.Post("/{id}/role", userHandler.UpdateRole)
			r.Post("/{id}/points", userHandler.UpdatePoints)
			r.Post("/{id}/ban", userHandler.Ban)
			r.Post("/{id}/unban", userHandler.Unban)
		})
	})
}
