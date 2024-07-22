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
		logger,
		testService,
	}, nil
}

func (h *TestHandler) IndexHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v, err := h.testService.HelloWorld()
			if err != nil {
				h.logger.Error(err.Error())
				http.Error(w, constants.InternalServerErrorMsg, http.StatusInternalServerError)
			}

			if err := helpers.Encode(
				w,
				http.StatusOK,
				"test endpoint successfully running",
				v,
			); err != nil {
				h.logger.Error(err.Error())
				http.Error(w, constants.InternalServerErrorMsg, http.StatusInternalServerError)
			}
		},
	)
}

// interface
