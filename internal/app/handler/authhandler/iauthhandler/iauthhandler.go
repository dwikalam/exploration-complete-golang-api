package iauthhandler

import "net/http"

type AuthHandler interface {
	HandleRegister() http.Handler
	HandleLogin() http.Handler
}
