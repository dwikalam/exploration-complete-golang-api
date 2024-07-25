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

type TestHandler struct {
	logger      interfaces.Logger
	testService *services.TestService
}

func NewTestHandler(
	logger interfaces.Logger,
	testService *services.TestService,
) (TestHandler, error) {
	if logger == nil || testService == nil {
		return TestHandler{}, errors.New("logger or testService is nil")
	}

	return TestHandler{
		logger:      logger,
		testService: testService,
	}, nil
}

func (h *TestHandler) HandleHelloWorldResponse() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
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
		},
	)
}

func (h *TestHandler) HandleTimeoutExceededResponse() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
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
		},
	)
}
