// Package db gestiona la conexion al motor de base de datos MySQL/MariaDB
// y configura el pool de conexiones del servidor.
package db

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

// NewMySQL abre una conexion al servidor MySQL/MariaDB usando variables de entorno
// y espera activamente hasta que la base de datos este disponible (util en Docker Compose).
//
// Variables de entorno requeridas:
//   - DB_HOST, DB_PORT, DB_NAME, DB_USER, DB_PASSWORD
//
// Variables de entorno opcionales (con valores por defecto):
//   - DB_MAX_OPEN_CONNS   (defecto 25): conexiones abiertas maximas en el pool.
//   - DB_MAX_IDLE_CONNS   (defecto 10): conexiones inactivas maximas en el pool.
//   - DB_CONN_MAX_LIFETIME (defecto 5m): tiempo de vida maximo de una conexion.
//
// Retorna error si no se puede establecer conexion en 30 segundos.
func NewMySQL() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		user, password, host, port, name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Configuracion del pool de conexiones.
	db.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 10))
	db.SetConnMaxLifetime(getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))

	// Espera activa: necesaria cuando la DB inicia en paralelo (Docker Compose).
	if err := waitForDB(db); err != nil {
		return nil, err
	}

	slog.Info("database connection established", "host", host, "name", name)

	return db, nil
}

// waitForDB reintenta la conexion cada 2 segundos hasta que la DB responda
// o se agote un timeout de 30 segundos.
//
// Entrada:  db, instancia sql.DB ya abierta.
// Salida:   nil si la DB responde a Ping(), error de timeout en caso contrario.
func waitForDB(db *sql.DB) error {
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(2 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout waiting for database connection")
		case <-ticker.C:
			if err := db.Ping(); err == nil {
				return nil
			}
			slog.Info("waiting for database...")
		}
	}
}
