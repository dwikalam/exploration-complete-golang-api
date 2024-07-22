package repositories

import (
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type TestRepository struct {
	logger interfaces.Logger
	db     interfaces.Database
}

func NewTestRepo(
	logger interfaces.Logger,
	db interfaces.Database,
) (TestRepository, error) {
	if db == nil {
		return TestRepository{}, errors.New("error interfaces.Database is nil")
	}

	return TestRepository{
		logger,
		db,
	}, nil
}

func (repo *TestRepository) GetAllTest() (any, error) {
	rows, err := repo.db.Access().Query("SELECT * FROM test")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return nil, nil
}
