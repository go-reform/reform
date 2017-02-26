package internal

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/reform.v1"
)

// Logger is our custom logger with Debugf method.
type Logger struct {
	printf reform.Printf
	debug  bool
}

// NewLogger creates a new logger.
func NewLogger(prefix string, debug bool) *Logger {
	var flags int
	if debug {
		flags = log.Ldate | log.Lmicroseconds | log.Lshortfile
	}

	l := log.New(os.Stderr, prefix, flags)
	return &Logger{
		printf: func(format string, args ...interface{}) {
			l.Output(3, fmt.Sprintf(format, args...))
		},
		debug: debug,
	}
}

// Debugf prints message only when Logger debug flag is set to true.
func (l *Logger) Debugf(format string, args ...interface{}) {
	if l.debug {
		l.printf(format, args...)
	}
}

// Printf calls l.Output to print to the logger.
// Arguments are handled in the manner of fmt.Printf.
func (l *Logger) Printf(format string, args ...interface{}) {
	l.printf(format, args...)
}

// Fatalf is equivalent to l.Printf() followed by a call to os.Exit(1) (or panic for debug logger).
func (l *Logger) Fatalf(format string, args ...interface{}) {
	l.printf(format, args...)
	if l.debug {
		// panic instead of os.Exit(1) to see output (SQL queries, failed assertions, etc.) in tests
		panic(fmt.Sprintf(format, args...))
	}
	os.Exit(1)
}
