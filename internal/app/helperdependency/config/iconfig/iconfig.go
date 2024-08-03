package iconfig

import "time"

type Getter interface {
	ServerConfigGetter
	DbConfigGetter
}

type ServerConfigGetter interface {
	ServerAddressConfigGetter
	ServerTimeoutConfigGetter
}

type ServerAddressConfigGetter interface {
	GetServerHost() string
	GetServerPort() int
}

type ServerTimeoutConfigGetter interface {
	GetServerTimeoutMessage() string
	GetServerReadTimeout() time.Duration
	GetServerWriteTimeout() time.Duration
	GetServerIdleTimeout() time.Duration
	GetServerHandlerTimeout() string
}

type DbConfigGetter interface {
	DbPsqlConfigGetter
}

type DbPsqlConfigGetter interface {
	GetDbPsqlDSN() string
	GetDbPsqlDriver() string
}
