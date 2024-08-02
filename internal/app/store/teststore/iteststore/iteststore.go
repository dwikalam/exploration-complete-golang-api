package iteststore

import "context"

type TestStorer interface {
	GetAll(ctx context.Context) (any, error)
	SimpleQuery(ctx context.Context) (int, error)
}
