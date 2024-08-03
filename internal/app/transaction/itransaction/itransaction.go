package itransaction

import "context"

type Manager interface {
	Run(
		ctx context.Context,
		callback func(ctx context.Context) error,
	) error
}
