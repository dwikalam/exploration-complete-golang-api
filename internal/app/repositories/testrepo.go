package repositories

import (
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type TestRepository struct {
	logger interfaces.Logger
	db     interfaces.DbManager
}

func NewTestRepo(logger interfaces.Logger, db interfaces.DbManager) (TestRepository, error) {
	if logger == nil || db == nil {
		return TestRepository{}, errors.New("logger or db is nil")
	}

	return TestRepository{
		logger: logger,
		db:     db,
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
