package derp

import (
	"time"
)

var reporters []Reporter

// Connect is called at system startup, and adds a new reporter to the slice of reporters to be notified when we report an error
func Connect(reporter Reporter) {
	reporters = append(reporters, reporter)
}

// New generates a new Error object
func New(location string, message string, innerError error, details ...interface{}) *Error {

	result := Error{
		Location:  location,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().Unix(),
	}

	// If we have an InnerError to work with...
	if result.InnerError != nil {

		// set the InnerError
		result.InnerError = Inspector(innerError)
		result.Code = result.InnerError.Code
	}

	// If we still don't have an HTTP error code, then default to 500.
	if result.Code == 0 {
		result.Code = 500
	}

	return &result
}

// NewWithCode generates a new Error object
func NewWithCode(location string, message string, innerError error, code int, details ...interface{}) *Error {

	return &Error{
		Location:   location,
		Code:       code,
		Message:    message,
		Details:    details,
		TimeStamp:  time.Now().Unix(),
		InnerError: Inspector(innerError),
	}
}
