package irepo

import (
	"context"

	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/repodto"
)

type UserRepository interface {
	GetByEmail(ctx context.Context, email string) (repodto.UserRet, error)
}
