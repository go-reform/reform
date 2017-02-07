package internal

import (
	"fmt"
	"log"
	"os"
)

// Logger is our custom logger with Debugf method.
type Logger struct {
	logger *log.Logger
	debug  bool
}

// NewLogger creates a new logger.
func NewLogger(prefix string, debug bool) *Logger {
	var flags int
	if debug {
		flags = log.Lshortfile
	}
	return &Logger{
		logger: log.New(os.Stderr, prefix, flags),
		debug:  debug,
	}
}

// Debugf prints message only when Logger debug flag is set to true.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.debug {
		l.logger.Output(2, fmt.Sprintf(format, args...))
	}
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.logger.Output(2, fmt.Sprintf(format, args...))
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1).
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.logger.Output(2, fmt.Sprintf(format, args...))
	os.Exit(1)
}
