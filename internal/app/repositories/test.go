package repositories

import (
	"context"
	"errors"
	"fmt"

	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type Test struct {
	db interfaces.DbQuerier
}

func NewTest(
	db interfaces.DbQuerier,
) (Test, error) {
	if db == nil {
		return Test{}, errors.New("nil db")
	}

	return Test{
		db: db,
	}, nil
}

func (r *Test) GetAll(ctx context.Context) (any, error) {
	const (
		query string = `
			SELECT 
				*
			FROM
				test
		`
	)

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return nil, nil
}

func (r *Test) SimpleQuery(ctx context.Context) (int, error) {
	var result int

	err := r.db.QueryRowContext(ctx, "SELECT 1").Scan(&result)

	if err != nil {
		return 0, fmt.Errorf("simple query failed: %v", err)
	}

	return result, nil
}
