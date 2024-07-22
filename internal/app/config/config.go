package config

import "github.com/dwikalam/ecommerce-service/internal/app/helpers"

type Config struct {
	serverConfig
	dbConfig
}

type serverConfig struct {
	ServerHost string
	ServerPort string
}

type dbConfig struct {
	PsqlURL string
}

func NewConfig() (Config, error) {
	if err := helpers.LoadEnv(); err != nil {
		return Config{}, err
	}

	serverConfig, err := newServerConfig()
	if err != nil {
		return Config{}, err
	}

	dbConfig, err := newDbConfig()
	if err != nil {
		return Config{}, err
	}

	return Config{
			serverConfig: serverConfig,
			dbConfig:     dbConfig,
		},
		nil
}

func newServerConfig() (serverConfig, error) {
	return serverConfig{
			ServerHost: helpers.GetEnvValue("SERVER_HOST", "localhost"),
			ServerPort: helpers.GetEnvValue("SERVER_PORT", "5000"),
		},
		nil
}

func newDbConfig() (dbConfig, error) {
	return dbConfig{
			PsqlURL: helpers.GetEnvValue("POSTGRESQL_URL", ""),
		},
		nil
}
