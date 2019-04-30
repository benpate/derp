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
func Wrap(inner *Error, location string, message string, details ...interface{}) *Error {

	result := Error{
		Location:   location,
		InnerError: inner,
		Message:    message,
		Details:    details,
		TimeStamp:  time.Now().Truncate(1 * time.Second),
		Code:       CodeInternalError,
	}

	if inner != nil {
		result.Code = inner.Code
	}

	return &result
}
