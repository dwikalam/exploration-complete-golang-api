package routes

import (
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/handlers"
	"github.com/dwikalam/ecommerce-service/internal/app/types/interfaces"
)

func NewHttpHandler(
	logger interfaces.Logger,
	testHandler *handlers.Test,
	authHandler *handlers.Auth,
) http.Handler {
	mux := http.NewServeMux()

	if testHandler != nil {
		mux.Handle("GET /api/v1/test", testHandler.HandleHelloWorldResponse())
		mux.Handle("GET /api/v1/test/timeout", testHandler.HandleTimeoutExceededResponse())
		mux.Handle("GET /api/v1/test/transaction", testHandler.HandleTimeoutExceededResponse())
	}

	if authHandler != nil {
		mux.Handle("POST /api/v1/register", authHandler.HandleRegister())
		mux.Handle("POST /api/v1/login", authHandler.HandleLogin())
	}

	return mux
}
