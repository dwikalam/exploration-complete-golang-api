package config

import (
	"errors"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type EnvConfig struct {
	server serverEnvConfig
	db     dbEnvConfig
}

func NewEnvConfig() (EnvConfig, error) {
	var (
		serverEnvConfig serverEnvConfig
		dbEnvConfig     dbEnvConfig

		err error
	)

	if err := godotenv.Load(); err != nil {
		return EnvConfig{}, err
	}

	serverEnvConfig, err = newServerEnvConfig()
	if err != nil {
		return EnvConfig{}, nil
	}

	dbEnvConfig, err = newDbEnvConfig()
	if err != nil {
		return EnvConfig{}, nil
	}

	return EnvConfig{
		server: serverEnvConfig,
		db:     dbEnvConfig,
	}, nil
}

func (c *EnvConfig) GetServerHost() string {
	return c.server.host
}

func (c *EnvConfig) GetServerPort() int {
	return c.server.port
}

func (c *EnvConfig) GetServerTimeoutMessage() string {
	return c.server.timeoutMessage
}

func (c *EnvConfig) GetServerReadTimeout() time.Duration {
	return c.server.readTimeout
}

func (c *EnvConfig) GetServerWriteTimeout() time.Duration {
	return c.server.writeTimeout
}

func (c *EnvConfig) GetServerIdleTimeout() time.Duration {
	return c.server.idleTimeout
}

func (c *EnvConfig) GetServerHandlerTimeout() time.Duration {
	return c.server.handlerTimeout
}

func (c *EnvConfig) GetDbPsqlDSN() string {
	return c.db.psqlConfig.dsn
}

func (c *EnvConfig) GetDbPsqlDriver() string {
	return c.db.psqlConfig.driver
}

type serverEnvConfig struct {
	host string
	port int

	timeoutMessage string
	readTimeout    time.Duration
	writeTimeout   time.Duration
	idleTimeout    time.Duration
	handlerTimeout time.Duration
}

func newServerEnvConfig() (serverEnvConfig, error) {
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
		return serverEnvConfig{}, err
	}

	readTimeout, err = time.ParseDuration(getEnvValue("SERVER_READ_TIMEOUT_MS", "1000") + "ms")
	if err != nil {
		return serverEnvConfig{}, err
	}
	writeTimeout, err = time.ParseDuration(getEnvValue("SERVER_WRITE_TIMEOUT_MS", "2000") + "ms")
	if err != nil {
		return serverEnvConfig{}, err
	}
	idleTimeout, err = time.ParseDuration(getEnvValue("SERVER_IDLE_TIMEOUT_MS", "60000") + "ms")
	if err != nil {
		return serverEnvConfig{}, err
	}
	handlerTimeout, err = time.ParseDuration(getEnvValue("SERVER_HANDLER_TIMEOUT_MS", "1000") + "ms")
	if err != nil {
		return serverEnvConfig{}, err
	}
	timeoutMessage = getEnvValue("SERVER_TIMEOUT_MESSAGE", "Timeout!")

	return serverEnvConfig{
		host: host,
		port: port,

		timeoutMessage: timeoutMessage,
		readTimeout:    readTimeout,
		writeTimeout:   writeTimeout,
		idleTimeout:    idleTimeout,
		handlerTimeout: handlerTimeout,
	}, nil
}

type dbEnvConfig struct {
	psqlConfig psqlConfig
}

func newDbEnvConfig() (dbEnvConfig, error) {
	psqlConfig, err := newPsqlConfig()
	if err != nil {
		return dbEnvConfig{}, errors.New("creating psqlConfig failed")
	}

	return dbEnvConfig{
		psqlConfig: psqlConfig,
	}, nil
}

type psqlConfig struct {
	driver string
	dsn    string
}

func newPsqlConfig() (psqlConfig, error) {
	return psqlConfig{
		driver: getEnvValue("POSTGRESQL_DRIVER", "pgx"),
		dsn:    getEnvValue("POSTGRESQL_DSN", ""),
	}, nil
}

func getEnvValue(envKey string, defaultValue string) string {
	if v, ok := os.LookupEnv(envKey); ok {
		return v
	}

	return defaultValue
}
