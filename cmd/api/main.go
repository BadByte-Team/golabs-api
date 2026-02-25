// Package main es el punto de entrada del servidor GoLabs API.
// Configura el logger estructurado, carga la configuracion, establece la conexion
// a la base de datos, construye el router HTTP y gestiona el apagado graceful.
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/joho/godotenv"

	"golabs-api/internal/infrastructure/config"
	"golabs-api/internal/infrastructure/db"
	http_if "golabs-api/internal/interfaces/http"
)

func main() {
	// Logger estructurado en JSON para integracion con Loki, Datadog, CloudWatch, etc.
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

	// Cargar configuracion estatica desde YAML.
	cfg, err := config.Load("configs/config.yaml")
	if err != nil {
		slog.Error("error loading config", "error", err)
		os.Exit(1)
	}

	// Cargar variables de entorno desde .env si existe. En produccion se usan
	// las variables del sistema directamente sin archivo .env.
	if err := godotenv.Load(); err != nil {
		slog.Warn("no .env file found, using system environment variables")
	}

	// Establecer conexion a la base de datos. La funcion espera activamente
	// hasta 30s si la DB aun esta iniciando (util en Docker Compose).
	database, err := db.NewMySQL()
	if err != nil {
		slog.Error("database connection failed", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	// Construir el router con todos los modulos registrados.
	router := http_if.NewRouter(database)

	// Configuracion del servidor HTTP con timeouts para prevenir resource exhaustion.
	server := &http.Server{
		Addr:         ":" + strconv.Itoa(cfg.Server.Port),
		Handler:      router,
		ReadTimeout:  5 * time.Second,  // tiempo maximo para leer el request completo
		WriteTimeout: 10 * time.Second, // tiempo maximo para escribir la respuesta
		IdleTimeout:  30 * time.Second, // tiempo maximo de inactividad de conexion keep-alive
	}

	// Iniciar el servidor en una goroutine separada para que el hilo principal
	// pueda esperar la senal de apagado.
	go func() {
		slog.Info("server starting", "app", cfg.App.Name, "port", cfg.Server.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			slog.Error("server error", "error", err)
			os.Exit(1)
		}
	}()

	shutdown(server)
}

// shutdown espera una senal SIGTERM o SIGINT y ejecuta el apagado graceful del servidor.
// Da hasta 15 segundos para que las peticiones en vuelo terminen antes de forzar el cierre.
//
// Entrada:  server, instancia http.Server en ejecucion.
func shutdown(server *http.Server) {
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop
	slog.Info("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		slog.Warn("forced shutdown", "error", err)
	}

	slog.Info("server stopped")
}
