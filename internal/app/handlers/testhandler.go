package handlers

import (
	"errors"
	"net/http"

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
	if testService == nil {
		return TestHandler{}, errors.New("error *services.TestService is nil")
	}

	return TestHandler{
		logger:      logger,
		testService: testService,
	}, nil
}

func (h *TestHandler) HandleResponseHelloWorld() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v, err := h.testService.HelloWorld()
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
				"test endpoint successfully running",
				v,
			)
		},
	)
}

// interface
