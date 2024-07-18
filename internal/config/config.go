package config

import "os"

var (
	HostServer string
	PortServer string
)

func init() {
	HostServer = getConfigValue("SERVER_HOST", "localhost")
	PortServer = getConfigValue("SERVER_PORT", "5000")
}

func getConfigValue(envKey string, defaultValue string) string {
	v, ok := os.LookupEnv(envKey)

	if !ok {
		return defaultValue
	}

	return v
}
