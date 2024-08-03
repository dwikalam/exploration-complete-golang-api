package testsvc

import (
	"context"
	"errors"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/store/teststore/iteststore"
	"github.com/dwikalam/ecommerce-service/internal/app/transaction/itransaction"
)

type Test struct {
	txManager itransaction.Manager
	testStore iteststore.Storer
}

func NewTest(
	txManager itransaction.Manager,
	testStore iteststore.Storer,
) (Test, error) {
	if txManager == nil {
		return Test{}, errors.New("nil txManager")
	}

	if testStore == nil {
		return Test{}, errors.New("nil testStore")
	}

	return Test{
		txManager: txManager,
		testStore: testStore,
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
		simpleTransaction = func(ctx context.Context) error {
			_, err := s.testStore.SimpleQuery(ctx)

			return err
		}

		result = func() string {
			err := s.txManager.Run(ctx, simpleTransaction)
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
