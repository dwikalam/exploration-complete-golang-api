package customerr

type DatabaseAlreadyConnectedError struct {
}

func (e *DatabaseAlreadyConnectedError) Error() string {
	return "error database is already connected"
}
