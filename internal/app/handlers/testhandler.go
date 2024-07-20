package handlers

import (
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/constants"
	"github.com/dwikalam/ecommerce-service/internal/app/helpers"
	"github.com/dwikalam/ecommerce-service/internal/app/services"
	"github.com/dwikalam/ecommerce-service/internal/app/types"
)

type TestHandler struct {
	testService *services.TestService
}

func NewTestHandler() *TestHandler {
	return &TestHandler{
		testService: services.GetTestServiceInstance(),
	}
}

func (h *TestHandler) IndexHandler() http.Handler {
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			v, err := h.testService.HelloWorld()
			if err != nil {
				http.Error(w, constants.InternalServerErrorMsg, http.StatusInternalServerError)
			}

			response := types.Response[string]{
				Message: "test endpoint successfully running",
				Data:    v,
			}
			if err := helpers.Encode(w, http.StatusOK, response); err != nil {
				http.Error(w, constants.InternalServerErrorMsg, http.StatusInternalServerError)
			}
		},
	)
}
