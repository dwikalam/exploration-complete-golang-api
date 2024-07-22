package interfaces

import (
	"context"
	"database/sql"
)

type Database interface {
	Access() *sql.DB
	Health(ctx context.Context) (map[string]string, error)
	Disconnect() error
}
