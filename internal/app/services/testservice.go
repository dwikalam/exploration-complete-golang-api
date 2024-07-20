package services

import (
	"github.com/dwikalam/ecommerce-service/internal/app/repositories"
)

type TestService struct {
	testRepo *repositories.TestRepository
}

var (
	instance *TestService
)

func GetTestServiceInstance() *TestService {
	if instance != nil {
		return instance
	}

	return &TestService{
		testRepo: repositories.GetTestRepositoryInstance(),
	}
}

func (s *TestService) HelloWorld() (string, error) {
	s.testRepo.GetAllTest()

	const v = "Hello, World!"

	return v, nil
}
