package icommon

import (
	"context"

	"github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"
)

type Validator interface {
	Valid(ctx context.Context) wrappertype.ProblemsMap
}
