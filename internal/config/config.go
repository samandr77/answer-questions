package config

import (
	"fmt"
	"log"
	"os"
	"strconv"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type ServerConfig struct {
	HTTPPort        string
	ReadTimeout     int
	WriteTimeout    int
	IdleTimeout     int
	ShutdownTimeout int
	RequestTimeout  int
}

func Load() *Config {
	return &Config{
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "password"),
			Name:     getEnv("DB_NAME", "questions_db"),
		},
		Server: ServerConfig{
			HTTPPort:        getEnv("HTTP_PORT", "8080"),
			ReadTimeout:     getEnvInt("HTTP_READ_TIMEOUT", 15),
			WriteTimeout:    getEnvInt("HTTP_WRITE_TIMEOUT", 15),
			IdleTimeout:     getEnvInt("HTTP_IDLE_TIMEOUT", 60),
			ShutdownTimeout: getEnvInt("HTTP_SHUTDOWN_TIMEOUT", 30),
			RequestTimeout:  getEnvInt("HTTP_REQUEST_TIMEOUT", 5),
		},
	}
}

func (db *DatabaseConfig) DSN() string {
	return fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		db.Host, db.Port, db.User, db.Password, db.Name)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	return value
}

func getEnvInt(key string, defaultValue int) int {
	value := os.Getenv(key)
	if value == "" {
		return defaultValue
	}
	intVal, err := strconv.Atoi(value)
	if err != nil {
		log.Printf("Ошибка преобразования %s в число: %v, используем значение по умолчанию %d", key, err, defaultValue)
		return defaultValue
	}
	return intVal
}
