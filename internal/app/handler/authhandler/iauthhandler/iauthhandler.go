package iauthhandler

import "net/http"

type Handler interface {
	HandleRegister() http.Handler
	HandleLogin() http.Handler
}
