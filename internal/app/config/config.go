package config

import "os"

var (
	ServerHost = getConfigValue("SERVER_HOST", "localhost")
	ServerPort = getConfigValue("SERVER_PORT", "5000")
)

var (
	PsqlURL = getConfigValue("POSTGRESQL_URL", "")
)

func getConfigValue(envKey string, defaultValue string) string {
	v, ok := os.LookupEnv(envKey)

	if !ok {
		return defaultValue
	}

	return v
}
