package isqldb

import (
	"context"
	"database/sql"
)

type Manager interface {
	Querier
	HealthChecker
	Disconnector
}

type Querier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type Accessor interface {
	Access() *sql.DB
}

type HealthChecker interface {
	CheckHealth(ctx context.Context) (map[string]string, error)
}

type Disconnector interface {
	Disconnect() error
}
