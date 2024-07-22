package interfaces

import (
	"context"
	"database/sql"
)

type DBAccessor interface {
	Access() *sql.DB
}

type DBHealthChecker interface {
	Health(ctx context.Context) (map[string]string, error)
}

type DBDisconnector interface {
	Disconnect() error
}
