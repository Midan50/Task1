package config

import (
	"log"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Host      string `yaml:"host" env-required:"true"`
	Port      string `yaml:"port" env-required:"true"`
	DBUrl     string `yaml:"db_url" env-required:"true"`
	JWTSecret string `yaml:"jwt_secret" env-required:"true"`
}

func MustLoad() *Config {
	var cfg Config
	if err := cleanenv.ReadConfig(".env.yaml", &cfg); err != nil {
		log.Fatalf("Failed to load config: %s", err)
	}
	return &cfg
}
