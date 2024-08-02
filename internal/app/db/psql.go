package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/transaction"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Psql struct {
	db *sql.DB
}

func NewPsql(psqlURL string) (Psql, error) {
	db, err := sql.Open("pgx", psqlURL)
	if err != nil {
		return Psql{}, err
	}

	return Psql{
		db: db,
	}, nil
}

func (p *Psql) Access() *sql.DB {
	return p.db
}

func (p *Psql) QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row {
	if tx, ok := transaction.ExtractTxFromContext(ctx); ok {
		return tx.QueryRowContext(ctx, query, args...)
	}

	return p.db.QueryRowContext(ctx, query, args...)
}

func (p *Psql) QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error) {
	if tx, ok := transaction.ExtractTxFromContext(ctx); ok {
		return tx.QueryContext(ctx, query, args...)
	}

	return p.db.QueryContext(ctx, query, args...)
}

func (p *Psql) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	if tx, ok := transaction.ExtractTxFromContext(ctx); ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return p.db.ExecContext(ctx, query, args...)
}

func (p *Psql) PrepareContext(ctx context.Context, query string) (*sql.Stmt, error) {
	if tx, ok := transaction.ExtractTxFromContext(ctx); ok {
		return tx.PrepareContext(ctx, query)
	}

	return p.db.PrepareContext(ctx, query)
}

func (p *Psql) CheckHealth(ctx context.Context) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(ctx, time.Second*1)
	defer cancel()

	err := p.db.PingContext(ctx)
	if err != nil {
		stats := map[string]string{
			"status": "down",
			"error":  fmt.Sprintf("error db down: %v", err),
		}

		return stats, err
	}

	dbStats := p.db.Stats()

	stats := map[string]string{
		"status":              "up",
		"message":             "It's healthy",
		"open_connections":    strconv.Itoa(dbStats.OpenConnections),
		"in_use":              strconv.Itoa(dbStats.InUse),
		"idle":                strconv.Itoa(dbStats.Idle),
		"wait_count":          strconv.Itoa(int(dbStats.WaitCount)),
		"wait_duration":       strconv.Itoa(int(dbStats.WaitDuration)),
		"max_idle_closed":     strconv.Itoa(int(dbStats.MaxIdleClosed)),
		"max_lifetime_closed": strconv.Itoa(int(dbStats.MaxLifetimeClosed)),
	}

	if dbStats.OpenConnections > 40 {
		stats["message"] = "The database is experiencing heavy load."
	}

	if dbStats.WaitCount > 1000 {
		stats["message"] = "The database has a high number of wait events, indicating potential bottlenecks."
	}

	if dbStats.MaxIdleClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many idle connections are being closed, consider revising the connection pool settings."
	}

	if dbStats.MaxLifetimeClosed > int64(dbStats.OpenConnections)/2 {
		stats["message"] = "Many connections are being closed due to max lifetime, consider increasing max lifetime or revising the connection usage pattern."
	}

	return stats, nil
}

func (p *Psql) Disconnect() error {
	var dbName string

	if err := p.db.QueryRow("SELECT current_database()").Scan(&dbName); err != nil {
		log.Printf("error fetching database name: %v", err)
	}

	if err := p.db.Close(); err != nil {
		return err
	}

	return nil
}
