package iauthsvc

import (
	"context"

	"github.com/dwikalam/ecommerce-service/internal/app/service/authsvc/authsvcdto"
)

type AuthServicer interface {
	RegisterUser(
		ctx context.Context,
		fullName string,
		email string,
		password string,
	) (authsvcdto.RegisteredUser, error)

	ValidateLoginAttempt(
		ctx context.Context,
		email string,
		plainPassword string,
	) error
}
