package configs

import (
	"log"
	"os"
)

type Config struct {
	SERVER_PORT string
	DB          DBConfig
	Auth        AuthConfig
}

type DBConfig struct {
	DB_HOST     string
	DB_PORT     string
	DB_NAME     string
	DB_USER     string
	DB_PASSWORD string
}

type AuthConfig struct {
	Secret string
}

func LoadConfig() *Config {
	// err := godotenv.Load()
	// if err != nil {
	// 	log.Println("[CFG] Error loading .env fole, using default config")
	// }
	if os.Getenv("SERVER_PORT") == "" || os.Getenv("DB_HOST") == "" || os.Getenv("DB_PORT") == "" || os.Getenv("DB_NAME") == "" || os.Getenv("DB_USER") == "" || os.Getenv("DB_PASSWORD") == "" || os.Getenv("SECRET") == "" {
		log.Println("[CFG] Some environment variables are not set, using default config")
	}

	return &Config{
		SERVER_PORT: os.Getenv("SERVER_PORT"),
		DB: DBConfig{
			DB_HOST:     os.Getenv("DB_HOST"),
			DB_PORT:     os.Getenv("DB_PORT"),
			DB_NAME:     os.Getenv("DB_NAME"),
			DB_USER:     os.Getenv("DB_USER"),
			DB_PASSWORD: os.Getenv("DB_PASSWORD"),
		},
		Auth: AuthConfig{
			Secret: os.Getenv("SECRET"),
		},
	}
}
