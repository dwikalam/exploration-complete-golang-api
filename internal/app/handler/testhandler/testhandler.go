package testhandler

import (
	"errors"
	"net/http"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/constant"
	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/codec/defaultcodec"
	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/logger/ilogger"
	"github.com/dwikalam/ecommerce-service/internal/app/service/testsvc/itestsvc"
)

type Test struct {
	logger      ilogger.Logger
	testService itestsvc.TestServicer
}

func NewTest(
	logger ilogger.Logger,
	testService itestsvc.TestServicer,
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

			defaultcodec.Encode[any](
				w,
				http.StatusInternalServerError,
				constant.InternalServerErrorMsg,
				nil,
			)

			return
		}

		defaultcodec.Encode(
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

			defaultcodec.Encode[any](
				w,
				http.StatusInternalServerError,
				constant.InternalServerErrorMsg,
				nil,
			)

			return
		}

		defaultcodec.Encode[any](
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

			defaultcodec.Encode[any](
				w,
				http.StatusInternalServerError,
				constant.InternalServerErrorMsg,
				nil,
			)

			return
		}

		defaultcodec.Encode[any](
			w,
			http.StatusOK,
			"Success",
			v,
		)
	})
}
