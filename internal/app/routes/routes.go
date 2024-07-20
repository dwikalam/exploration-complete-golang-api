package routes

import (
	"net/http"

	"github.com/dwikalam/ecommerce-service/internal/app/handlers"
)

func NewHttpHandler(
	testHandler *handlers.TestHandler,
) http.Handler {
	mux := http.NewServeMux()

	if testHandler != nil {
		mux.Handle("GET /api/v1/test", testHandler.IndexHandler())
	}

	return mux
}
