package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server serverConfig
	Db     dbConfig
}

func New() (Config, error) {
	var (
		serverConfig serverConfig
		dbConfig     dbConfig

		err error
	)

	if err := godotenv.Load(); err != nil {
		return Config{}, err
	}

	serverConfig, err = newServerConfig()
	if err != nil {
		return Config{}, nil
	}

	dbConfig, err = newDbConfig()
	if err != nil {
		return Config{}, nil
	}

	return Config{
		Server: serverConfig,
		Db:     dbConfig,
	}, nil
}

type serverConfig struct {
	Host string
	Port int

	TimeoutMessage string
	ReadTimeout    time.Duration
	WriteTimeout   time.Duration
	IdleTimeout    time.Duration
	HandlerTimeout time.Duration
}

func newServerConfig() (serverConfig, error) {
	var (
		host string
		port int

		timeoutMessage string
		readTimeout    time.Duration
		writeTimeout   time.Duration
		idleTimeout    time.Duration
		handlerTimeout time.Duration

		err error
	)

	host = getEnvValue("SERVER_HOST", "localhost")
	port, err = strconv.Atoi(getEnvValue("SERVER_PORT", "8080"))
	if err != nil {
		return serverConfig{}, err
	}

	readTimeout, err = time.ParseDuration(getEnvValue("SERVER_READ_TIMEOUT_MS", "1000") + "ms")
	if err != nil {
		return serverConfig{}, err
	}
	writeTimeout, err = time.ParseDuration(getEnvValue("SERVER_WRITE_TIMEOUT_MS", "2000") + "ms")
	if err != nil {
		return serverConfig{}, err
	}
	idleTimeout, err = time.ParseDuration(getEnvValue("SERVER_IDLE_TIMEOUT_MS", "60000") + "ms")
	if err != nil {
		return serverConfig{}, err
	}
	handlerTimeout, err = time.ParseDuration(getEnvValue("SERVER_HANDLER_TIMEOUT_MS", "1000") + "ms")
	if err != nil {
		return serverConfig{}, err
	}
	timeoutMessage = getEnvValue("SERVER_TIMEOUT_MESSAGE", "Timeout!")

	return serverConfig{
		Host: host,
		Port: port,

		TimeoutMessage: timeoutMessage,
		ReadTimeout:    readTimeout,
		WriteTimeout:   writeTimeout,
		IdleTimeout:    idleTimeout,
		HandlerTimeout: handlerTimeout,
	}, nil
}

type dbConfig struct {
	PsqlURL string
}

func newDbConfig() (dbConfig, error) {
	return dbConfig{
		PsqlURL: getEnvValue("POSTGRESQL_URL", ""),
	}, nil
}

func getEnvValue(envKey string, defaultValue string) string {
	if v, ok := os.LookupEnv(envKey); ok {
		return v
	}

	return defaultValue
}
