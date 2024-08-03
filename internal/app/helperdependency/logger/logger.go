package logger

import (
	"io"
	"log"
)

type Logger struct {
	infoLogger  *log.Logger
	warnLogger  *log.Logger
	errorLogger *log.Logger
}

func New(stdout, stderr io.Writer) Logger {
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

	return Logger{
		infoLogger:  newLogger(stdout, infoLoggerPrefix),
		warnLogger:  newLogger(stdout, warnLoggerPrefix),
		errorLogger: newLogger(stderr, errorLoggerPrefix),
	}
}

func (l *Logger) Info(message string) {
	l.infoLogger.Output(2, message)
}

func (l *Logger) Warn(message string) {
	l.warnLogger.Output(2, message)
}

func (l *Logger) Error(message string) {
	l.errorLogger.Output(2, message)
}

func newLogger(out io.Writer, prefix string) *log.Logger {
	return log.New(out, prefix, log.LstdFlags|log.Lshortfile)
}
