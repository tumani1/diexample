// Package errors provide implementations of custom errors for the echo framework.
package errors

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// HTTPError represents an error that occurred while handling a request.
type HTTPError struct {
	Code        int
	Cause       error
	Message     string
	Description interface{}
}

// New creates a new HTTPError instance.
func New(code int, message ...string) *HTTPError {
	var e = &HTTPError{
		Code:    code,
		Message: http.StatusText(code),
	}

	if len(message) > 0 {
		e.Message = strings.Join(message, " ")
	}

	return e
}

// Wrap creates a new HTTPError instance with wrap origin error.
func Wrap(code int, err error, message ...string) *echo.HTTPError {
	var e = &echo.HTTPError{
		Code:     code,
		Internal: err,
	}

	if len(message) > 0 {
		e.Message = strings.Join(message, " ")
	}

	return e
}

// Error makes it compatible with `error` interface.
func (e *HTTPError) Error() string {
	var msg = e.Message
	if e.Cause != nil {
		if len(e.Message) > 0 {
			msg = e.Message + ": " + e.Cause.Error()
		} else {
			msg = e.Cause.Error()
		}
	}

	return fmt.Sprintf("code=%d, message=%s", e.Code, msg)
}
