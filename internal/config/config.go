package config

import (
	"fmt"
	"os"
	"time"
)

// Helper pour récupérer les variables d'environnement avec valeurs par défaut
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}

// Configuration par défaut
const (
	defaultDBName     = "music_share"
	defaultDBUser     = "postgres"
	defaultDBPassword = "Drogobeats1995*"
	defaultDBHost     = "localhost"
	defaultDBPort     = "5432"
	defaultSSLMode    = "disable"
	defaultJWTSecret  = "your-super-secret-key"
	defaultServerPort = "10000"
)

// Durée de token fixe
const TokenExpiration = 24 * time.Hour

type Config struct {
	DB struct {
		Host     string
		User     string
		Password string
		Name     string
		Port     string
		SSLMode  string
	}
	JWT struct {
		Secret     string
		Expiration time.Duration
	}
	Server struct {
		Port string
	}
}

func LoadConfig() *Config {
	// Ajouter ces logs de debug
	fmt.Printf("DB_HOST from env: %s\n", os.Getenv("DB_HOST"))
	fmt.Printf("DB_PORT from env: %s\n", os.Getenv("DB_PORT"))
	fmt.Printf("DB_USER from env: %s\n", os.Getenv("DB_USER"))

	return &Config{
		DB: struct {
			Host     string
			User     string
			Password string
			Name     string
			Port     string
			SSLMode  string
		}{
			Host:     getEnv("DB_HOST", defaultDBHost),
			User:     getEnv("DB_USER", defaultDBUser),
			Password: getEnv("DB_PASSWORD", defaultDBPassword),
			Name:     getEnv("DB_NAME", defaultDBName),
			Port:     getEnv("DB_PORT", defaultDBPort),
			SSLMode:  getEnv("DB_SSL_MODE", defaultSSLMode),
		},
		JWT: struct {
			Secret     string
			Expiration time.Duration
		}{
			Secret:     getEnv("JWT_SECRET", defaultJWTSecret),
			Expiration: TokenExpiration,
		},
		Server: struct {
			Port string
		}{
			Port: getEnv("SERVER_PORT", defaultServerPort),
		},
	}
}

func (c *Config) GetDBDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.DB.Host,
		c.DB.User,
		c.DB.Password,
		c.DB.Name,
		c.DB.Port,
		c.DB.SSLMode,
	)
}
