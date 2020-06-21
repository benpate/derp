package derp

import (
	"time"
)

// New returns a new Error object
func New(code int, location string, message string, details ...interface{}) *Error {

	result := Error{
		Location:  location,
		Code:      code,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().Truncate(1 * time.Second),
	}

	return &result
}

// Wrap encapsulates an existing derp.Error
func Wrap(inner error, location string, message string, details ...interface{}) *Error {

	result := Error{
		Location:  location,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().Truncate(1 * time.Second),
		Code:      CodeInternalError,
	}

	// If we're wrapping another derp, then bubble its values up.
	if innerDerp, ok := inner.(*Error); ok {
		result.InnerError = innerDerp
		result.Code = innerDerp.Code
	} else if inner != nil {
		// Otherwise, append generic error to the end of the "details" section.
		result.Details = append(result.Details, inner)
	}

	return &result
}
