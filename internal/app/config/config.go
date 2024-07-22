package config

import "os"

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
			ServerHost: getConfigValue("SERVER_HOST", "localhost"),
			ServerPort: getConfigValue("SERVER_PORT", "5000"),
		},
		nil
}

func newDbConfig() (dbConfig, error) {
	return dbConfig{
			PsqlURL: getConfigValue("POSTGRESQL_URL", ""),
		},
		nil
}

func getConfigValue(envKey string, defaultValue string) string {
	v, ok := os.LookupEnv(envKey)

	if !ok {
		return defaultValue
	}

	return v
}
