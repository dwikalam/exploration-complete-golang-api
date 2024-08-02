package response

type Default[T any] struct {
	Message string `json:"message"`
	Data    T      `json:"data"`
}
