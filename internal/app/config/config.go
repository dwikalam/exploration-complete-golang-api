package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	serverConfig
	dbConfig
}

func New() (Config, error) {
	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	return Config{
		serverConfig: newServerConfig(),
		dbConfig:     newDbConfig(),
	}, nil
}

type serverConfig struct {
	ServerHost string
	ServerPort string
}

func newServerConfig() serverConfig {
	return serverConfig{
		ServerHost: getEnvValue("SERVER_HOST", "localhost"),
		ServerPort: getEnvValue("SERVER_PORT", "5000"),
	}
}

type dbConfig struct {
	PsqlURL string
}

func newDbConfig() dbConfig {
	return dbConfig{
		PsqlURL: getEnvValue("POSTGRESQL_URL", ""),
	}
}

func getEnvValue(envKey string, defaultValue string) string {
	v, ok := os.LookupEnv(envKey)

	if !ok {
		return defaultValue
	}

	return v
}
