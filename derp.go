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
func New(location string, code int, message string, innerError error, details ...interface{}) *Error {

	result := Error{
		Location:  location,
		Code:      code,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().Unix(),
	}

	if innerError != nil {

		switch e := innerError.(type) {
		case *Error:

			// Embed the innerError into the new object we're creating.
			result.InnerError = e

			if result.Code == 0 {
				result.Code = e.Code
			}

		default:

			// Other, unrecognized kinds of errors get wrapped in a derp.Error, so that we can embed them correctly.
			result.InnerError = &Error{
				Location: "Embedded Error",
				Message:  e.Error(),
				Details:  []interface{}{e},
				Code:     CodeInternalError,
			}

			result.Code = CodeInternalError
		}
	}

	// If we still don't have an HTTP error code, then default to CodeInternalError.
	if result.Code == 0 {
		result.Code = CodeInternalError
	}

	return &result
}
