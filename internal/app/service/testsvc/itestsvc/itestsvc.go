package itestsvc

import (
	"context"
	"time"
)

type Servicer interface {
	HelloWorld(ctx context.Context) (string, error)
	OperateFor(ctx context.Context, d time.Duration) error
	Transaction(ctx context.Context) (string, error)
}
