package interfaces

import (
	"context"
	"database/sql"
)

type DbManager interface {
	DbQuerier
	DbHealthChecker
	DbDisconnector
}

type DbQuerier interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	PrepareContext(ctx context.Context, query string) (*sql.Stmt, error)
}

type DbAccessor interface {
	Access() *sql.DB
}

type DbHealthChecker interface {
	CheckHealth(ctx context.Context) (map[string]string, error)
}

type DbDisconnector interface {
	Disconnect() error
}

type TransactionManager interface {
	Run(
		ctx context.Context,
		callback func(ctx context.Context) error,
	) error
}
