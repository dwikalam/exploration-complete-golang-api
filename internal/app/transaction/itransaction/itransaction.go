package itransaction

import "context"

type TransactionManager interface {
	Run(
		ctx context.Context,
		callback func(ctx context.Context) error,
	) error
}
