// Package config carga la configuracion estatica del servidor desde un archivo YAML.
// Las variables sensibles (contrasenas, secretos) deben usarse via variables de entorno,
// no en el archivo de configuracion.
package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

// Config contiene la configuracion de alto nivel del servidor.
// Se deserializa desde configs/config.yaml al arrancar la aplicacion.
type Config struct {
	// Server agrupa los parametros del servidor HTTP.
	Server struct {
		Port int `yaml:"port"` // puerto en el que escucha el servidor
	} `yaml:"server"`

	// App agrupa metadatos de la aplicacion usados en logging y monitoreo.
	App struct {
		Name string `yaml:"name"` // nombre del servicio (aparece en logs estructurados)
		Env  string `yaml:"env"`  // entorno: "development", "staging", "production"
	} `yaml:"app"`
}

// Load lee el archivo YAML en path y lo deserializa en un Config.
//
// Entrada:  path, ruta relativa o absoluta al archivo de configuracion YAML.
// Salida:   puntero a Config o error si el archivo no existe o tiene formato invalido.
func Load(path string) (*Config, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var cfg Config
	if err := yaml.Unmarshal(file, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
