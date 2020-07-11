package derp

import "time"

// Error represents a runtime error.  It includes
type Error struct {
	Code       int           `json:"code"`       // Numeric error code (such as an HTTP status code) to report to the client.
	Location   string        `json:"location"`   // Function name (or other location description) of where the error occurred
	Message    string        `json:"message"`    // Primary (top-level) error message for this error
	TimeStamp  time.Time     `json:"timestamp"`  // Unix Epoch timestamp of the date/time when this error was created
	Details    []interface{} `json:"details"`    // Additional information related to this error message, such as parameters to the function that caused the error.
	InnerError *Error        `json:"innerError"` // An underlying error object used to identify the root cause of this error.
}

// Error implements the Error interface, which allows derp.Error objects to be
// used anywhere a standard error is used.
func (err *Error) Error() string {
	return err.Location + ": " + err.Message
}

// ErrorCode returns the error Code embedded in this Error.  This is useful for matching
// interfaces in other package.
func (err *Error) ErrorCode() int {
	return err.Code
}

// Unwrap supports Go 1.13+ error unwrapping
func (err *Error) Unwrap() error {

	// If we have an InnerError, then return that
	if err.InnerError != nil {
		return err.InnerError
	}

	// Otherise, look for a generic error wrapped into the details
	if length := len(err.Details); length > 0 {

		// If the last item is an error, then return that
		if err, ok := err.Details[length-1].(error); ok {
			return err
		}
	}

	// Fall through means that we couldn't find anything.
	return nil
}

// RootCause digs into the error stack and returns the original error
// that caused the DERP
func (err *Error) RootCause() *Error {

	if (err.InnerError != nil) && (err.InnerError.Message != "") {
		return err.InnerError.RootCause()
	}

	return err
}

// NotFound returns TRUE if the error `Code` is a 404 / Not Found error.
func (err *Error) NotFound() bool {
	return err.Code == CodeNotFoundError
}
