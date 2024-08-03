package teststore

import (
	"context"
	"errors"
	"fmt"

	"github.com/dwikalam/ecommerce-service/internal/app/db/sqldb/isqldb"
)

type SQLStore struct {
	sqldb isqldb.Querier
}

func NewSQLStore(
	sqldb isqldb.Querier,
) (SQLStore, error) {
	if sqldb == nil {
		return SQLStore{}, errors.New("nil sqldb")
	}

	return SQLStore{
		sqldb: sqldb,
	}, nil
}

func (store *SQLStore) GetAll(ctx context.Context) (any, error) {
	const (
		query string = `
			SELECT 
				*
			FROM
				test
		`
	)

	rows, err := store.sqldb.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("QueryContext failed: %v", err)
	}

	defer rows.Close()

	return nil, nil
}

func (store *SQLStore) SimpleQuery(ctx context.Context) (int, error) {
	var result int

	if err := store.sqldb.QueryRowContext(ctx, "SELECT 1").Scan(&result); err != nil {
		return 0, fmt.Errorf("QueryRowContext failed: %v", err)
	}

	return result, nil
}
