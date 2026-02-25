// Package userhttp implementa los handlers HTTP del modulo de usuarios
// y el registro de sus rutas en el router principal.
package userhttp

import (
	"database/sql"

	"github.com/go-chi/chi/v5"

	"golabs-api/internal/infrastructure/security"
	accessmw "golabs-api/internal/interfaces/http/middleware/access"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/middleware/ratelimit"
	refreshtokenapp "golabs-api/internal/refreshtoken/application"
	refreshtokendomain "golabs-api/internal/refreshtoken/domain"
	refreshtokeninfra "golabs-api/internal/refreshtoken/infrastructure"
	userapp "golabs-api/internal/user/application"
	userdomain "golabs-api/internal/user/domain"
	userinfra "golabs-api/internal/user/infrastructure"
)

// RegisterRoutes registra todas las rutas de usuario y autenticacion en el router dado.
//
// Rutas publicas (con rate limiting de login):
//   - POST /auth/login      autenticar usuario, retorna par de tokens
//   - POST /auth/register   crear cuenta nueva
//   - POST /auth/refresh    rotar refresh token (token rotation)
//   - POST /auth/logout     revocar refresh token (idempotente)
//
// Rutas protegidas (JWT + LoadUser + no baneado + rate limit por usuario):
//   - GET  /users/                    lista de usuarios (admin)
//   - POST /users/                    crear usuario (admin)
//   - GET  /users/search?q=           busqueda parcial por username
//   - GET  /users/by-username/{name}  busqueda exacta por username
//   - GET  /users/{id}                detalle de usuario
//   - POST /users/{id}/update         actualizar perfil (propio o admin)
//   - POST /users/{id}/password       cambiar contrasena (propio o admin)
//
// Rutas de admin exclusivas (JWT + rol admin):
//   - POST /admin/users/{id}/role     cambiar rol
//   - POST /admin/users/{id}/points   actualizar puntos
//   - POST /admin/users/{id}/ban      banear usuario
//   - POST /admin/users/{id}/unban    desbanear usuario
func RegisterRoutes(r chi.Router, db *sql.DB, jwtSvc *security.JWTService) {
	repo := userinfra.NewUserRepository(db)
	rtRepo := refreshtokeninfra.New(db)

	// Use cases de refresh token.
	issueRT := refreshtokenapp.NewIssueRefreshTokenUseCase(rtRepo)
	refreshRT := refreshtokenapp.NewRefreshAccessTokenUseCase(rtRepo, repo, jwtSvc, issueRT)
	revokeRT := refreshtokenapp.NewRevokeRefreshTokenUseCase(rtRepo)

	// Use cases de usuario.
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

	// Handlers.
	authHandler := NewAuthHandler(loginUC, createUserUC, issueRT, refreshRT, revokeRT)
	userHandler := NewUserHandler(
		createUserUC, getUserByIDUC, getUserByUsernameUC, searchByUsernameUC,
		listUsersUC,
		updateUserUC, changePasswordUC, updateRoleUC, updatePointsUC,
		banUserUC, unbanUserUC,
	)

	// Rutas publicas con rate limiting de autenticacion.
	r.Route("/auth", func(r chi.Router) {
		r.Use(ratelimit.LoginRateLimit)
		r.Post("/login", authHandler.Login)
		r.Post("/register", authHandler.Register)
		r.Post("/refresh", authHandler.Refresh)
		r.Post("/logout", authHandler.Logout) // siempre 204, nunca falla para el cliente
	})

	// Rutas protegidas: requieren JWT valido + usuario no baneado.
	r.Group(func(r chi.Router) {
		r.Use(authmw.JWTAuth(jwtSvc))
		r.Use(authmw.LoadUser(repo)) // recarga rol y flag banned desde BD
		r.Use(accessmw.RequireNotBanned)
		r.Use(ratelimit.UserRateLimit)

		r.Route("/users", func(r chi.Router) {
			r.With(accessmw.RequireRole(userdomain.RoleAdmin)).Get("/", userHandler.List)
			r.With(accessmw.RequireRole(userdomain.RoleAdmin)).Post("/", userHandler.Create)
			r.Get("/search", userHandler.Search)
			r.Get("/by-username/{username}", userHandler.GetByUsername)
			r.Get("/{id}", userHandler.GetByID)
			r.With(accessmw.RequireSelfOrAdmin).Post("/{id}/update", userHandler.Update)
			r.With(accessmw.RequireSelfOrAdmin).Post("/{id}/password", userHandler.ChangePassword)
		})

		// Rutas de administracion exclusivas para el rol admin.
		r.Route("/admin/users", func(r chi.Router) {
			r.Use(accessmw.RequireRole(userdomain.RoleAdmin))
			r.Post("/{id}/role", userHandler.UpdateRole)
			r.Post("/{id}/points", userHandler.UpdatePoints)
			r.Post("/{id}/ban", userHandler.Ban)
			r.Post("/{id}/unban", userHandler.Unban)
		})
	})
}

// Validacion en tiempo de compilacion: asegura que MySQLRefreshTokenRepository
// implementa la interfaz RefreshTokenRepository. Evita errores silenciosos.
var _ refreshtokendomain.RefreshTokenRepository = (*refreshtokeninfra.MySQLRefreshTokenRepository)(nil)
