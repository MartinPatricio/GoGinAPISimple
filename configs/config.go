package configs

import (
	"errors"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type ConfigDB struct {
	// Servidor
	Port    string
	GinMode string

	// Base de datos
	Database DB

	// Autenticación
	APISecretKey string
	JWTSecret    string

	// Aplicación
	AppEnv   string
	LogLevel string
}

type DB struct {
	User     string
	Password string
	Host     string
	Port     int
	NameDB   string
	SSLMode  string
}

func Load() *ConfigDB {
	if err := godotenv.Load(); err != nil {
		panic("Tenemos un problema con las variables de enterno o no esta el archivo env")
	}

	config := &ConfigDB{
		Port:    getEnv("PORT", "8080"),
		GinMode: getEnv("GIN_MODE", "debug"),
		Database: DB{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnvAsInt("DB_PORT", 5432),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", ""),
			NameDB:   getEnv("DB_NAME", "user_api_db"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		APISecretKey: getEnv("API_SECRET_KEY", ""),
		JWTSecret:    getEnv("JWT_SECRET", ""),
		AppEnv:       getEnv("APP_ENV", "development"),
		LogLevel:     getEnv("LOG_LEVEL", "info"),
	}
	if err := config.Validate(); err != nil {
		log.Fatal("Configuration error:", err)
	}
	return config
}

// VALIDATE STATUS DATABASE

func (conn *ConfigDB) Validate() error {
	if conn.Database.Password == "" {
		return errors.New("DB_PASSWORD i required")
	}
	if conn.APISecretKey == "" {
		return errors.New("API_SECRET_KEY i required")
	}

	if len(conn.APISecretKey) < 32 {
		return errors.New("API_SECRET_KEY must key be at least 32 characters long")
	}
	return nil
}

// HELPERS FOR SECRETKEY

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
		log.Printf("Invalid integer value for %s: %s, using default: %d", key, value, defaultValue)
	}
	return defaultValue
}
