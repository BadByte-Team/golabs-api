// Package http es el punto de entrada de la capa de interfaces HTTP.
// NewRouter construye el chi.Mux raiz, registra el middleware global y
// delega el registro de rutas especificas a cada modulo de la aplicacion.
package http

import (
	"database/sql"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"

	"golabs-api/internal/health"
	"golabs-api/internal/infrastructure/security"
	authmw "golabs-api/internal/interfaces/http/middleware/auth"
	"golabs-api/internal/interfaces/http/middleware/bodylimit"

	challengehttp "golabs-api/internal/challenges/interfaces"
	eventhttp "golabs-api/internal/event/interfaces"
	eventteamhttp "golabs-api/internal/eventteam/interfaces"
	userhttp "golabs-api/internal/user/interfaces"
)

// maxBodyBytes es el limite de tamano del body de las peticiones HTTP (1 MiB).
// Previene ataques de agotamiento de memoria con payloads gigantes.
const maxBodyBytes = 1 << 20

// NewRouter construye y retorna el router principal de la API.
//
// Orden del middleware global:
//  1. RequestID:   asigna un ID unico a cada peticion para correlacion en logs.
//  2. RealIP:      extrae la IP real del cliente desde headers de proxy (X-Forwarded-For).
//  3. Recoverer:   captura panics y retorna HTTP 500 en lugar de crashear el proceso.
//  4. ErrorLogger: registra errores 5xx con slog.
//  5. MaxBodySize: rechaza bodies mayores a 1 MiB con HTTP 413.
//  6. CORS:        configura encabezados de cross-origin segun ALLOWED_ORIGINS.
//
// Estructura de rutas:
//   - /healthz/live  y  /healthz/ready: probes de orquestador (sin versioning).
//   - /health:                          path legacy para compatibilidad.
//   - /api/v1/**:                       rutas versionadas de todos los modulos.
func NewRouter(db *sql.DB) *chi.Mux {
	r := chi.NewRouter()

	// Middleware global aplicado a todas las rutas.
	r.Use(chimw.RequestID)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(authmw.ErrorLogger)
	r.Use(bodylimit.MaxBodySize(maxBodyBytes))
	r.Use(corsMiddleware())

	// Health checks (sin prefijo de version — los orquestadores esperan paths fijos).
	h := health.NewHandler(db)
	r.Get("/healthz/live", h.Live)
	r.Get("/healthz/ready", h.Ready)
	r.Get("/health", h.ServeHTTP) // path legacy

	jwtSvc, err := security.NewJWTService()
	if err != nil {
		panic("JWT no configurado: " + err.Error())
	}

	// Rutas versionadas: cada modulo registra su propio grupo de rutas.
	r.Route("/api/v1", func(r chi.Router) {
		userhttp.RegisterRoutes(r, db, jwtSvc)
		eventhttp.RegisterRoutes(r, db, jwtSvc)
		eventteamhttp.RegisterRoutes(r, db, jwtSvc)
		challengehttp.RegisterRoutes(r, db, jwtSvc)
	})

	return r
}

// corsMiddleware construye el handler CORS a partir de la variable de entorno ALLOWED_ORIGINS.
//
// Configuracion:
//   - ALLOWED_ORIGINS="http://localhost:3000,https://app.example.com"
//   - Si ALLOWED_ORIGINS no esta definida, se permite cualquier origen ("*") para desarrollo.
//   - AllowCredentials se activa solo cuando se especifican origenes concretos (requerimiento del estandar CORS).
func corsMiddleware() func(http.Handler) http.Handler {
	originsEnv := os.Getenv("ALLOWED_ORIGINS")
	var origins []string
	if originsEnv == "" {
		origins = []string{"*"}
	} else {
		for _, o := range strings.Split(originsEnv, ",") {
			if s := strings.TrimSpace(o); s != "" {
				origins = append(origins, s)
			}
		}
	}

	return cors.Handler(cors.Options{
		AllowedOrigins:   origins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-Request-Id"},
		ExposedHeaders:   []string{"X-Request-Id"},
		AllowCredentials: originsEnv != "",
		MaxAge:           300,
	})
}
