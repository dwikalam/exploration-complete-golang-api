package authhandler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/handler/authhandler/dto/authreqdto"
	"github.com/dwikalam/ecommerce-service/internal/app/handler/authhandler/dto/authrespdto"
	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/codec/codec"
	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/logger/ilogger"
	"github.com/dwikalam/ecommerce-service/internal/app/service/authsvc/authsvcdto"
	"github.com/dwikalam/ecommerce-service/internal/app/service/authsvc/iauthsvc"
	"github.com/dwikalam/ecommerce-service/internal/app/type/wrappertype"
)

type Auth struct {
	logger      ilogger.Logger
	authService iauthsvc.Servicer
}

func New(
	logger ilogger.Logger,
	authService iauthsvc.Servicer,
) (Auth, error) {
	if logger == nil {
		return Auth{}, errors.New("nil logger")
	}

	if authService == nil {
		return Auth{}, errors.New("nil authService")
	}

	return Auth{
		logger:      logger,
		authService: authService,
	}, nil
}

func (h *Auth) HandleRegister() http.Handler {
	const (
		errMsg     string = "failed to register user"
		successMsg string = "user registered successfully"
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var (
			payloads *authreqdto.RegisterUser
			problems wrappertype.ProblemsMap
			err      error

			registerUserSvcDto authsvcdto.RegisteredUser

			registerHandlerDto authrespdto.RegisteredUser
		)

		payloads, problems, err = codec.DecodeValid[*authreqdto.RegisterUser](r)
		if problems != nil {
			h.logger.Error(fmt.Sprintf("decode valid failed: %s", problems))

			codec.Encode(
				w,
				http.StatusBadRequest,
				errMsg,
				problems,
			)

			return
		}
		if err != nil {
			const errData = "request json payload not valid"

			h.logger.Error(fmt.Sprintf("decode valid failed: %v", err))

			codec.Encode(
				w,
				http.StatusBadRequest,
				errMsg,
				errData,
			)

			return
		}

		registerUserSvcDto, err = h.authService.RegisterUser(
			r.Context(),
			payloads.FullName,
			payloads.Email,
			payloads.Password,
		)
		if err != nil {
			const errData = "email already exist"

			h.logger.Error(fmt.Sprintf("register user service failed: %v", err))

			codec.Encode(
				w,
				http.StatusBadRequest,
				errMsg,
				errData,
			)

			return
		}

		registerHandlerDto = authrespdto.RegisteredUser{
			ID: registerUserSvcDto.ID,
		}

		codec.Encode(
			w,
			http.StatusCreated,
			successMsg,
			registerHandlerDto,
		)
	})
}

func (h *Auth) HandleLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
