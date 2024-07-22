package services

import (
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/repositories"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type TestService struct {
	logger   interfaces.Logger
	testRepo *repositories.TestRepository
}

func NewTestService(
	logger interfaces.Logger,
	testRepo *repositories.TestRepository,
) (TestService, error) {
	if testRepo == nil {
		return TestService{}, errors.New("error * repositories.TestRepository is nil")
	}

	return TestService{
		logger,
		testRepo,
	}, nil
}

func (s *TestService) HelloWorld() (string, error) {
	s.testRepo.GetAllTest()

	const v = "Hello, World!"

	return v, nil
}
