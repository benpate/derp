package derp

import "time"

// Error represents a runtime error.  It includes
type Error struct {
	Location   string        `json:"location"`   // Function name (or other location description) of where the error occurred
	Message    string        `json:"message"`    // Primary (top-level) error message for this error
	Code       int           `json:"code"`       // Numeric error code (such as an HTTP status code) to report to the client.
	TimeStamp  time.Time     `json:"timestamp"`  // Unix Epoch timestamp of the date/time when this error was created
	Details    []interface{} `json:"details"`    // Additional information related to this error message, such as parameters to the function that caused the error.
	InnerError *Error        `json:"innerError"` // An underlying error object used to identify the root cause of this error.
}

// Error implements the Error interface, which allows derp.Error objects to be
// used anywhere a standard error is used.
func (err *Error) Error() string {
	return err.Location + ": " + err.Message
}

// RootCause digs into the error stack and returns the original error
// that caused the DERP
func (err *Error) RootCause() *Error {

	if (err.InnerError != nil) && (err.InnerError.Message != "") {
		return err.InnerError.RootCause()
	}

	return err
}

// Report sends this error to all configured plugins, to be reported via their various error reporting channels.
func (err *Error) Report() *Error {

	for _, plugin := range Plugins {
		plugin.Report(err)
	}

	return err
}

// NotFound returns TRUE if the error `Code` is a 404 / Not Found error.
func (err *Error) NotFound() bool {
	return err.Code == CodeNotFoundError
}
