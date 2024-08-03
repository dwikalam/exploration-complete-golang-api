package authreqdto

import (
	"context"
	"errors"

	"github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"
)

type LoginUser struct {
	Email    string
	Password string
}

func (p *LoginUser) Valid(ctx context.Context) wrappertype.ProblemsMap {
	var (
		problems = make(wrappertype.ProblemsMap)

		err error
	)

	if err = p.validatePayloadStructure(); err != nil {
		problems["payload"] = err.Error()

		return problems
	}

	return nil
}

func (p *LoginUser) validatePayloadStructure() error {
	if p.Email == "" || p.Password == "" {
		return errors.New("payload structure not valid")
	}

	return nil
}
