package void

import (
	"fmt"
	"log"
	"runtime"
)

type Error struct {
	s       string
	context interface{}
}

func (e *Error) Error() string {
	if e.context == nil {
		return e.s
	}

	return fmt.Sprintf("%s, '%v'", e.s, e.context)
}

func NewError(s string) error {
	return &Error{s: s}
}

func NewErrorContext(s string, context interface{}) error {
	return &Error{s: s, context: context}
}

// CatchPanic is used to catch any Panic and log exceptions to Stdout. It will also write the stack trace.
func CatchPanic(err *error, goRoutine string, functionName string) {
	if r := recover(); r != nil {
		// Capture the stack trace
		buf := make([]byte, 10000)
		runtime.Stack(buf, false)

		writeStdoutf(goRoutine, functionName, "PANIC Defered [%v] : Stack Trace : %v", r, string(buf))

		if err != nil {
			*err = fmt.Errorf("%v", r)
		}
	}
}

// writeStdout is used to write a system message directly to stdout.
func writeStdout(goRoutine string, functionName string, message string) {
	log.Printf("%s : %s : %s\n", goRoutine, functionName, message)
}

// writeStdoutf is used to write a formatted system message directly stdout.
func writeStdoutf(goRoutine string, functionName string, format string, a ...interface{}) {
	writeStdout(goRoutine, functionName, fmt.Sprintf(format, a...))
}
