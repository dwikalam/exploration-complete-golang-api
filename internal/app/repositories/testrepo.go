package repositories

import "github.com/dwikalam/ecommerce-service/internal/app/db"

type TestRepository struct {
	db *db.Database
}

var (
	instance *TestRepository
)

func GetTestRepositoryInstance() *TestRepository {
	if instance != nil {
		return instance
	}

	return &TestRepository{
		db: db.GetInstance(),
	}
}

func (repo *TestRepository) GetAllTest() (any, error) {
	rows, err := repo.db.Query("SELECT * FROM test")
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	return nil, nil
}
