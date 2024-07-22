package interfaces

import (
	"context"
	"database/sql"
)

type DbManager interface {
	dbAccessor
	dbHealthChecker
	dbDisconnector
}

type dbAccessor interface {
	Access() *sql.DB
}
type dbHealthChecker interface {
	CheckHealth(ctx context.Context) (map[string]string, error)
}

type dbDisconnector interface {
	Disconnect() error
}
