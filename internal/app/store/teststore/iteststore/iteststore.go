package iteststore

import "context"

type Storer interface {
	GetAll(ctx context.Context) (any, error)
	SimpleQuery(ctx context.Context) (int, error)
}
