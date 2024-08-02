package iuserstore

import (
	"context"

	"github.com/dwikalam/ecommerce-service/internal/app/store/userstore/userstoredto"
)

type UserStorer interface {
	GetByEmail(ctx context.Context, email string) (userstoredto.User, error)
	Create(ctx context.Context, fullName string, email string, password string) (userstoredto.User, error)
}
