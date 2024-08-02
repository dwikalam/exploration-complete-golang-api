package authsvcdto

import "github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"

type RegisteredUser struct {
	ID       wrappertype.DbID
	FullName string
	Email    string
	Password string
}
