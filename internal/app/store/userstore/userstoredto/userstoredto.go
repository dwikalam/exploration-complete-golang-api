package userstoredto

import (
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"
)

type User struct {
	ID        wrappertype.DbID `json:"id"`
	FullName  string           `json:"fullName"`
	Email     string           `json:"email"`
	Password  string           `json:"password"`
	CreatedAt time.Time        `json:"createdAt"`
	UpdatedAt time.Time        `json:"updatedAt"`
}
