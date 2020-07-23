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

	// If there is no error to wrap, then return nothing.
	// This makes some code simpler because it minimizes "if err... " checking.
	if inner == nil {
		return nil
	}

	result := Error{
		InnerError: inner,
		Location:   location,
		Message:    message,
		Details:    details,
		TimeStamp:  time.Now().Truncate(1 * time.Second),
		Code:       CodeInternalError,
	}

	// If we're wrapping another derp, then bubble its values up.
	if innerDerp, ok := inner.(*Error); ok {
		result.Code = innerDerp.Code
	} else {
		result.Details = append([]interface{}{inner.Error()}, result.Details...)
	}

	return &result
}
