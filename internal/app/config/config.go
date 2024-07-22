package config

import "github.com/dwikalam/ecommerce-service/internal/app/helpers"

type Config struct {
	serverConfig
	dbConfig
}

func New() (Config, error) {
	if err := helpers.LoadEnv(); err != nil {
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
		ServerHost: helpers.GetEnvValue("SERVER_HOST", "localhost"),
		ServerPort: helpers.GetEnvValue("SERVER_PORT", "5000"),
	}
}

type dbConfig struct {
	PsqlURL string
}

func newDbConfig() dbConfig {
	return dbConfig{
		PsqlURL: helpers.GetEnvValue("POSTGRESQL_URL", ""),
	}
}
