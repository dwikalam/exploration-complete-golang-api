package loggers

import (
	"io"
	"log"
)

type Default struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func NewDefault(stdout, stderr io.Writer) Default {
	const (
		resetColor  = "\033[0m"
		redColor    = "\033[31m"
		yellowColor = "\033[33m"
		greenColor  = "\033[32m"

		infoPlaceholder    = "INFO:    "
		warningPlaceholder = "WARNING: "
		errorPlaceholder   = "ERROR:   "
	)

	var (
		infoLoggerPrefix  = greenColor + infoPlaceholder + resetColor
		warnLoggerPrefix  = yellowColor + warningPlaceholder + resetColor
		errorLoggerPrefix = redColor + errorPlaceholder + resetColor
	)

	return Default{
		infoLogger:  newLogger(stdout, infoLoggerPrefix),
		warnLogger:  newLogger(stdout, warnLoggerPrefix),
		errorLogger: newLogger(stderr, errorLoggerPrefix),
	}
}

func (l *Default) Info(message string) {
	l.infoLogger.Output(2, message)
}

func (l *Default) Warn(message string) {
	l.warnLogger.Output(2, message)
}

func (l *Default) Error(message string) {
	l.errorLogger.Output(2, message)
}

func newLogger(out io.Writer, prefix string) *log.Logger {
	return log.New(out, prefix, log.LstdFlags|log.Lshortfile)
}
