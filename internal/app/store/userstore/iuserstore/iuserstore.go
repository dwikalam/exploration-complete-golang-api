package iuserstore

import (
	"context"

	"github.com/dwikalam/ecommerce-service/internal/app/store/userstore/userstoredto"
)

type Storer interface {
	GetByEmail(ctx context.Context, email string) (userstoredto.User, error)
	IsEmailExist(ctx context.Context, email string) (bool, error)
	Create(ctx context.Context, fullName string, email string, password string) (userstoredto.User, error)
}
