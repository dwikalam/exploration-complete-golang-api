package itesthandler

import "net/http"

type Handler interface {
	HandleHelloWorldResponse() http.Handler
	HandleTimeoutExceededResponse() http.Handler
	HandleTransactionTest() http.Handler
}
