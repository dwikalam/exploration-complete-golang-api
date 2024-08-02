package authrespdto

import "github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"

type RegisteredUser struct {
	ID wrappertype.DbID `json:"id"`
}
