package config

import "github.com/caarlos0/env/v11"

// Config agrupa la configuración del API leída desde variables de entorno.
type Config struct {
	DatabaseURL string `env:"DATABASE_URL,required"`
	Port        string `env:"PORT" envDefault:"8080"`
}

// Load lee la configuración desde el entorno.
func Load() (Config, error) {
	return env.ParseAs[Config]()
}
