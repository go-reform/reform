package internal

import (
	"fmt"
	"log"
	"os"
)

// Logger is our custom logger with Debugf method.
type Logger struct {
	*log.Logger
	debug bool
}

// NewLogger creates a new logger.
func NewLogger(prefix string, debug bool) *Logger {
	var flags int
	if debug {
		flags = log.Lshortfile
	}
	return &Logger{
		Logger: log.New(os.Stderr, prefix, flags),
		debug:  debug,
	}
}

// Debugf prints message only when Logger debug flag is set to true.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.debug {
		l.Output(2, fmt.Sprintf(format, args...))
	}
}
