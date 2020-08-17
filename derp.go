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

	if result.InnerError != nil {
		if innerDerp, ok := result.InnerError.(*Error); ok {
			result.Code = innerDerp.Code
		} else {
			result.Details = append([]interface{}{result.InnerError.Error()}, result.Details...)
		}
	}

	return &result
}
