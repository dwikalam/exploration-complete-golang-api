package db

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/config"
	"github.com/dwikalam/ecommerce-service/internal/app/types/customerr"
	_ "github.com/jackc/pgx/v5/stdlib"
)

type Database struct {
	*sql.DB
}

var (
	instance *Database
)

func Initialize() error {
	if instance != nil {
		return &customerr.DatabaseAlreadyConnectedError{}
	}

	db, err := sql.Open("pgx", config.PsqlURL)
	if err != nil {
		return err
	}

	instance = &Database{db}

	return nil
}

func GetInstance() *Database {
	return instance
}

func (db *Database) Health() (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*1)
	defer cancel()

	err := db.PingContext(ctx)
	if err != nil {
		stats := map[string]string{
			"status": "down",
			"error":  fmt.Sprintf("error db down: %v", err),
		}

		return stats, err
	}

	dbStats := db.Stats()

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

	if dbStats.OpenConnections > 40 { // Assuming 50 is the max for this example
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

func (db *Database) Disconnect() error {
	var dbName string

	if err := db.QueryRow("SELECT current_database()").Scan(&dbName); err != nil {
		log.Printf("error fetching database name: %v", err)
	}

	if err := db.Close(); err != nil {
		log.Printf("error disconnecting from database: %s", dbName)

		return err
	}

	fmt.Printf("Disconnected from database: %s", dbName)

	return nil
}
