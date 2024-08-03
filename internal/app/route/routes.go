package route

import (
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/handler/authhandler/iauthhandler"
	"github.com/dwikalam/ecommerce-service/internal/app/handler/testhandler/itesthandler"
)

func NewHttpHandler(
	testHandler itesthandler.Handler,
	authHandler iauthhandler.Handler,
) http.Handler {
	mux := http.NewServeMux()

	if testHandler != nil {
		mux.Handle("GET /api/v1/test", testHandler.HandleHelloWorldResponse())
		mux.Handle("GET /api/v1/test/timeout", testHandler.HandleTimeoutExceededResponse())
		mux.Handle("GET /api/v1/test/transaction", testHandler.HandleTransactionTest())
	}

	if authHandler != nil {
		mux.Handle("POST /api/v1/auth/register", authHandler.HandleRegister())
		mux.Handle("POST /api/v1/auth/login", authHandler.HandleLogin())
	}

	return mux
}
