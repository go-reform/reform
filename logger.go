package reform

import (
	"fmt"
	"strings"
	"time"
)

// Inspect returns suitable for logging representation of a query argument.
func Inspect(arg interface{}, addType bool) string {
	// do not merge cases, we want "arg == nil" to work
	var s string
	switch arg := arg.(type) {
	case string:
		s = fmt.Sprintf("%#q", arg)

	case *string:
		if arg == nil {
			s = "<nil>"
		} else {
			s = Inspect(*arg, false)
		}

	case fmt.Stringer:
		if arg == nil {
			s = "<nil>"
		} else {
			s = fmt.Sprintf("%s", arg)
		}

	case *int32:
		if arg == nil {
			s = "<nil>"
		} else {
			s = Inspect(*arg, false)
		}

	default:
		s = fmt.Sprintf("%v", arg)
	}

	if addType {
		s += fmt.Sprintf(" (%T)", arg)
	}
	return s
}

// Logger is responsible to log queries before and after their execution.
type Logger interface {
	// Before logs query before execution.
	Before(query string, args []interface{})

	// After logs query after execution.
	After(query string, args []interface{}, d time.Duration, err error)
}

// Printf is a (fmt.Printf|log.Printf|testing.T.Logf)-like function.
type Printf func(format string, a ...interface{})

// PrintfLogger is a simple query logger.
type PrintfLogger struct {
	LogTypes bool
	printf   Printf
}

// NewPrintfLogger creates a new simple query logger for any Printf-like function.
func NewPrintfLogger(printf Printf) *PrintfLogger {
	return &PrintfLogger{false, printf}
}

// Before logs query before execution.
func (pl *PrintfLogger) Before(query string, args []interface{}) {
	// fast path
	if args == nil {
		pl.printf(">>> %s", query)
		return
	}

	ss := make([]string, len(args))
	for i, arg := range args {
		ss[i] = Inspect(arg, pl.LogTypes)
	}

	pl.printf(">>> %s [%s]", query, strings.Join(ss, ", "))
}

// After logs query after execution.
func (pl *PrintfLogger) After(query string, args []interface{}, d time.Duration, err error) {
	// fast path
	if args == nil {
		msg := fmt.Sprintf("%s %s", query, d)
		if err != nil {
			msg += ": " + err.Error()
		}
		pl.printf("<<< %s", msg)
		return
	}

	ss := make([]string, len(args))
	for i, arg := range args {
		ss[i] = Inspect(arg, pl.LogTypes)
	}

	msg := fmt.Sprintf("%s [%s] %s", query, strings.Join(ss, ", "), d)
	if err != nil {
		msg += ": " + err.Error()
	}
	pl.printf("<<< %s", msg)
}

// check interface
var _ Logger = new(PrintfLogger)
