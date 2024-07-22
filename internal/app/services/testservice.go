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

func NewTestService(logger interfaces.Logger, testRepo *repositories.TestRepository) (TestService, error) {
	if logger == nil || testRepo == nil {
		return TestService{}, errors.New("logger or testRepo is nil")
	}

	return TestService{
		logger:   logger,
		testRepo: testRepo,
	}, nil
}

func (s *TestService) HelloWorld() (string, error) {
	const v = "Hello, World!"

	_, err := s.testRepo.GetAllTest()
	if err != nil {
		return "", err
	}

	return v, nil
}
