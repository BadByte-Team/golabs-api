package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func NewMySQL() (*sql.DB, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	name := os.Getenv("DB_NAME")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?parseTime=true&charset=utf8mb4&loc=UTC",
		user,
		password,
		host,
		port,
		name,
	)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return nil, err
	}

	// Pool de conexiones
	db.SetMaxOpenConns(getEnvInt("DB_MAX_OPEN_CONNS", 25))
	db.SetMaxIdleConns(getEnvInt("DB_MAX_IDLE_CONNS", 10))
	db.SetConnMaxLifetime(getEnvDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute))

	// Espera activa (DB en Docker)
	if err := waitForDB(db); err != nil {
		return nil, err
	}

	log.Println("[INFO] conexión a MariaDB establecida")

	return db, nil
}

func waitForDB(db *sql.DB) error {
	timeout := time.After(30 * time.Second)
	tick := time.Tick(2 * time.Second)

	for {
		select {
		case <-timeout:
			return fmt.Errorf("timeout esperando conexión a DB")
		case <-tick:
			if err := db.Ping(); err == nil {
				return nil
			}
			log.Println("[INFO] esperando base de datos...")
		}
	}
}
