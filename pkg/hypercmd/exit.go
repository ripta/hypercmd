package hypercmd

import (
	"errors"
	"fmt"
)

// ExitError is an error that carries the process exit code it should produce.
// Return one from a command's RunE, wrapping it with Exit, to request an exit code other than the default of 1.
type ExitError struct {
	Err  error
	Code int
}

// Exit wraps err so that ExitCode reports code for it. err may be nil.
func Exit(code int, err error) error {
	return &ExitError{Err: err, Code: code}
}

func (e *ExitError) Error() string {
	if e.Err == nil {
		return fmt.Sprintf("exit code %d", e.Code)
	}
	return e.Err.Error()
}

func (e *ExitError) Unwrap() error {
	return e.Err
}

// ExitCode returns the process exit code for err: 0 for a nil error, the
// code carried by err if it is (or wraps) an *ExitError, and 1 for any
// other non-nil error.
func ExitCode(err error) int {
	if err == nil {
		return 0
	}

	var ee *ExitError
	if errors.As(err, &ee) {
		return ee.Code
	}

	return 1
}
