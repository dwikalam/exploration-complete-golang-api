package interfaces

import (
	"context"

	"github.com/dwikalam/ecommerce-service/internal/app/types/customtype"
)

type Validator interface {
	Valid(ctx context.Context) customtype.ProblemsMap
}
