package main

import (
	"fmt"
	"log"
	"os"
)

// Logger is our custom logger with Debugf method.
type Logger struct {
	*log.Logger
	Debug bool
}

// NewLogger creates a new logger.
func NewLogger() *Logger {
	return &Logger{
		Logger: log.New(os.Stderr, "reform: ", 0),
		Debug:  false,
	}
}

// Debugf prints message only when Logger debug flag is set to true.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.Debug {
		l.Output(2, fmt.Sprintf(format, args...))
	}
}
