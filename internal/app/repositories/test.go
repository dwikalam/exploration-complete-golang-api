package repositories

import (
	"context"
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type Test struct {
	logger interfaces.Logger
	db     interfaces.DbQuerier
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

func (repo *Test) GetAll(ctx context.Context) (any, error) {
	const (
		query string = `
			SELECT 
				*
			FROM
				test
		`
	)

	rows, err := repo.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return nil, nil
}
