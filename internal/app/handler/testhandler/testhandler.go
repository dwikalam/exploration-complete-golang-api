package testhandler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/dwikalam/ecommerce-service/internal/app/constant"
	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/codec/codec"
	"github.com/dwikalam/ecommerce-service/internal/app/helperdependency/logger/ilogger"
	"github.com/dwikalam/ecommerce-service/internal/app/service/testsvc/itestsvc"
)

type Test struct {
	logger      ilogger.Logger
	testService itestsvc.Servicer
}

func New(
	logger ilogger.Logger,
	testService itestsvc.Servicer,
) (Test, error) {
	if logger == nil {
		return Test{}, errors.New("logger is nil")
	}

	if testService == nil {
		return Test{}, errors.New("testService is nil")
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
			h.logger.Error(fmt.Sprintf("HelloWorld service failed: %v", err))

			codec.Encode[any](
				w,
				http.StatusInternalServerError,
				constant.InternalServerErrorMsg,
				nil,
			)

			return
		}

		codec.Encode(
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
			h.logger.Error(fmt.Sprintf("OperateFor service failed: %v", err))

			codec.Encode[any](
				w,
				http.StatusInternalServerError,
				constant.InternalServerErrorMsg,
				nil,
			)

			return
		}

		codec.Encode[any](
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
			h.logger.Error(fmt.Sprintf("Transaction service failed: %v", err))

			codec.Encode[any](
				w,
				http.StatusInternalServerError,
				constant.InternalServerErrorMsg,
				nil,
			)

			return
		}

		codec.Encode[any](
			w,
			http.StatusOK,
			"Success",
			v,
		)
	})
}
