package itesthandler

import "net/http"

type TestHandler interface {
	HandleHelloWorldResponse() http.Handler
	HandleTimeoutExceededResponse() http.Handler
	HandleTransactionTest() http.Handler
}
