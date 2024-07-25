package services

import (
	"context"
	"errors"
	"time"

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

func (s *TestService) HelloWorld(ctx context.Context) (string, error) {
	vChannel := func() <-chan string {
		const v = "Hello, World!"

		ch := make(chan string)

		go func() {
			select {
			case <-ctx.Done():
			case ch <- v:
			}
		}()

		return ch
	}

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case v := <-vChannel():
		return v, nil
	}
}

func (s *TestService) OperateFor(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}
