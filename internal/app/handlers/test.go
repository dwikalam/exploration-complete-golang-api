package handlers

import (
	"errors"
	"net/http"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/constants"
	"github.com/dwikalam/ecommerce-service/internal/app/helpers"
	"github.com/dwikalam/ecommerce-service/internal/app/services"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

type Test struct {
	logger      interfaces.Logger
	testService *services.Test
}

func NewTest(
	logger interfaces.Logger,
	testService *services.Test,
) (Test, error) {
	if logger == nil || testService == nil {
		return Test{}, errors.New("logger or testService is nil")
	}

	return Test{
		logger:      logger,
		testService: testService,
	}, nil
}

func (h *Test) HandleHelloWorldResponse() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, err := h.testService.HelloWorld(r.Context())
		if err != nil {
			h.logger.Error(err.Error())

			helpers.Encode[any](
				w,
				http.StatusInternalServerError,
				constants.InternalServerErrorMsg,
				nil,
			)

			return
		}

		helpers.Encode(
			w,
			http.StatusOK,
			"HandleHelloWorldResponse successfully running",
			v,
		)
	})
}

func (h *Test) HandleTimeoutExceededResponse() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		const duration = time.Second * 5

		if err := h.testService.OperateFor(r.Context(), duration); err != nil {
			h.logger.Error(err.Error())

			helpers.Encode[any](
				w,
				http.StatusInternalServerError,
				constants.InternalServerErrorMsg,
				nil,
			)

			return
		}

		helpers.Encode[any](
			w,
			200,
			"HandleTimeoutExceededResponse successfully running",
			nil,
		)
	})
}

func (h *Test) HandleTransactionTest() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		v, err := h.testService.Transaction(r.Context())

		if err != nil {
			h.logger.Error(err.Error())

			helpers.Encode[any](
				w,
				http.StatusInternalServerError,
				constants.InternalServerErrorMsg,
				nil,
			)
		}

		helpers.Encode[any](
			w,
			http.StatusOK,
			"Success",
			v,
		)
	})
}
