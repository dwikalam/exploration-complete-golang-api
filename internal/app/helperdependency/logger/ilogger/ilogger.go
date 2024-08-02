package ilogger

type Logger interface {
	Informer
	Warner
	Error
}

type Informer interface {
	Info(message string)
}

type Warner interface {
	Warn(message string)
}

type Error interface {
	Error(message string)
}
