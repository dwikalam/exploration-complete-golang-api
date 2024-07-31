package services

import (
	"context"
	"errors"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/repositories"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type Test struct {
	logger    interfaces.Logger
	txManager interfaces.TransactionManager
	testRepo  *repositories.Test
}

func NewTest(
	logger interfaces.Logger,
	txManager interfaces.TransactionManager,
	testRepo *repositories.Test,
) (Test, error) {
	if logger == nil || txManager == nil || testRepo == nil {
		return Test{}, errors.New("there is nil dependency")
	}

	return Test{
		logger:    logger,
		txManager: txManager,
		testRepo:  testRepo,
	}, nil
}

func (s *Test) HelloWorld(ctx context.Context) (string, error) {
	var (
		vChannel = func() <-chan string {
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
	)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case v := <-vChannel():
		return v, nil
	}
}

func (s *Test) OperateFor(ctx context.Context, d time.Duration) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-time.After(d):
		return nil
	}
}

func (s *Test) Transaction(ctx context.Context) (string, error) {
	var (
		result = func() string {
			testTx := func(ctx context.Context) error {
				_, err := s.testRepo.GetAll(ctx)

				return err
			}

			err := s.txManager.Run(ctx, testTx)
			if err != nil {
				return "failed"
			}

			return "success"
		}

		vChannel = func() <-chan string {
			ch := make(chan string)

			go func() {
				defer close(ch)

				select {
				case <-ctx.Done():
				case ch <- result():
				}
			}()

			return ch
		}
	)

	select {
	case <-ctx.Done():
		return "", ctx.Err()
	case v := <-vChannel():
		return v, nil
	}
}
