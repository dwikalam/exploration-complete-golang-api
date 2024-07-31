package repositories

import (
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/repodto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type User struct {
	logger interfaces.Logger
	db     interfaces.DbAccessor
}

func NewUser(logger interfaces.Logger, db interfaces.DbAccessor) (User, error) {
	if logger == nil || db == nil {
		return User{}, errors.New("logger or db is nil")
	}

	return User{
		logger: logger,
		db:     db,
	}, nil
}

func (r *User) Register(repodto.UserRegister) error {
	r.db.Access().Begin()

	return nil
}
