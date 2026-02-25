// Package logging provee una funcion de construccion del logger estructurado
// de la aplicacion usando log/slog con salida en formato JSON.
package logging

import (
	"log/slog"
	"os"
)

// NewLogger crea un slog.Logger que escribe en logs/app.log en formato JSON
// y lo registra como logger por defecto del proceso.
//
// El directorio logs/ se crea automaticamente si no existe.
//
// Salida: puntero al logger configurado o error si no se puede crear/abrir el archivo.
//
// Nota: el archivo de log crece indefinidamente; en produccion se recomienda
// rotar con logrotate o delegar a stdout y capturar con el orchestrador (Docker/k8s).
func NewLogger() (*slog.Logger, error) {
	if err := os.MkdirAll("logs", 0755); err != nil {
		return nil, err
	}

	file, err := os.OpenFile(
		"logs/app.log",
		os.O_APPEND|os.O_CREATE|os.O_WRONLY,
		0644,
	)
	if err != nil {
		return nil, err
	}

	handler := slog.NewJSONHandler(file, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})

	logger := slog.New(handler)
	slog.SetDefault(logger)

	return logger, nil
}
