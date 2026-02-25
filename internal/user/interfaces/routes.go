package userhttp

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	"golabs-api/internal/infrastructure/security"
	accessmw "golabs-api/internal/interfaces/http/middleware/access"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/middleware/ratelimit"
	userapp "golabs-api/internal/user/application"
	userdomain "golabs-api/internal/user/domain"
	userinfra "golabs-api/internal/user/infrastructure"
)

// RegisterRoutes wires all user and auth routes into the given router.
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {
	repo := userinfra.NewUserRepository(db)

	// Use cases
	createUserUC := userapp.NewCreateUserUseCase(repo)
	getUserByIDUC := userapp.NewGetUserByIDUseCase(repo)
	getUserByUsernameUC := userapp.NewGetUserByUsernameUseCase(repo)
	searchByUsernameUC := userapp.NewSearchUserByUsernameUseCase(repo)
	listUsersUC := userapp.NewListUsersUseCase(repo)
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
		createUserUC, getUserByIDUC, getUserByUsernameUC, searchByUsernameUC,
		listUsersUC,
		updateUserUC, changePasswordUC, updateRoleUC, updatePointsUC,
		banUserUC, unbanUserUC,
	)

	// Public auth routes
	r.Route("/auth", func(r chi.Router) {
		r.Use(ratelimit.LoginRateLimit)
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
	})

	// Protected routes (authenticated + not banned)
	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuth(jwtSvc))
		r.Use(authmw.LoadUser(repo)) // loads Banned flag and fresh Role from DB
		r.Use(accessmw.RequireNotBanned)
		r.Use(ratelimit.UserRateLimit)

		r.Route("/users", func(r chi.Router) {
			// Admin: list all users (paginated)
			r.With(accessmw.RequireRole(userdomain.RoleAdmin)).Get("/", userHandler.List)

			// Admin: create user
			r.With(accessmw.RequireRole(userdomain.RoleAdmin)).Post("/", userHandler.Create)

			// Search: GET /users/search?q=<query>
			r.Get("/search", userHandler.Search)

			// Exact username lookup: GET /users/by-username/{username}
			r.Get("/by-username/{username}", userHandler.GetByUsername)

			// By ID
			r.Get("/{id}", userHandler.GetByID)
			r.With(accessmw.RequireSelfOrAdmin).Post("/{id}/update", userHandler.Update)
			r.With(accessmw.RequireSelfOrAdmin).Post("/{id}/password", userHandler.ChangePassword)
		})

		r.Route("/admin/users", func(r chi.Router) {
			r.Use(accessmw.RequireRole(userdomain.RoleAdmin))
			r.Post("/{id}/role", userHandler.UpdateRole)
			r.Post("/{id}/points", userHandler.UpdatePoints)
			r.Post("/{id}/ban", userHandler.Ban)
			r.Post("/{id}/unban", userHandler.Unban)
		})
	})
}
