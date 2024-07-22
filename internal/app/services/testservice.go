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
		return TestService{}, errors.New("*repositories.TestRepository is nil")
	}

	return TestService{
		logger,
		testRepo,
	}, nil
}

func (s *TestService) HelloWorld() (string, error) {
	_, err := s.testRepo.GetAllTest()
	if err != nil {
		return "", err
	}

	const v = "Hello, World!"

	return v, nil
}
