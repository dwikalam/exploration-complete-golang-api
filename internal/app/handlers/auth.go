package handlers

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/helpers"
	"github.com/dwikalam/ecommerce-service/internal/app/services"
	"github.com/dwikalam/ecommerce-service/internal/app/types/customtype"
	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/reqdto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/respdto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/dto/svcdto"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type Auth struct {
	logger      interfaces.Logger
	authService *services.Auth
}

func NewAuth(logger interfaces.Logger, authService *services.Auth) (Auth, error) {
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
			payloads *reqdto.RegisterUser
			problems customtype.ProblemsMap
			err      error

			registerSvcArg svcdto.RegisterUserArg
			registerSvcRet svcdto.RegisteredUserRet

			response respdto.RegisteredUser
		)

		payloads, problems, err = helpers.DecodeValid[*reqdto.RegisterUser](r)
		if problems != nil {
			h.logger.Error(fmt.Sprintf("%s", problems))

			helpers.Encode(
				w,
				http.StatusBadRequest,
				errMsg,
				problems,
			)

			return
		}
		if err != nil {
			const errData = "request json payload not valid"

			h.logger.Error(err.Error())

			helpers.Encode(
				w,
				http.StatusBadRequest,
				errMsg,
				errData,
			)

			return
		}

		registerSvcArg = svcdto.RegisterUserArg{
			FullName: payloads.FullName,
			Email:    payloads.Email,
			Password: payloads.Password,
		}
		registerSvcRet, err = h.authService.RegisterUser(r.Context(), &registerSvcArg)
		if err != nil {
			const errData = "email has been registered"

			h.logger.Error(err.Error())

			helpers.Encode(
				w,
				http.StatusBadRequest,
				errMsg,
				errData,
			)

			return
		}

		response = respdto.RegisteredUser{
			ID: registerSvcRet.ID,
		}

		helpers.Encode(
			w,
			http.StatusCreated,
			successMsg,
			response,
		)

	})
}

func (h *Auth) HandleLogin() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

	})
}
