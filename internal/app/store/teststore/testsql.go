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

func NewTest(
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
		return nil, err
	}

	defer rows.Close()

	return nil, nil
}

func (store *SQLStore) SimpleQuery(ctx context.Context) (int, error) {
	var result int

	err := store.sqldb.QueryRowContext(ctx, "SELECT 1").Scan(&result)

	if err != nil {
		return 0, fmt.Errorf("simple query failed: %v", err)
	}

	return result, nil
}
