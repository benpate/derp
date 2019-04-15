package derp

// Error represents a runtime error.  It includes
type Error struct {
	Location   string        // Function name (or other location description) of where the error occurred
	Code       int           // Numeric error code (such as an HTTP status code) to report to the client.
	Message    string        // Primary (top-level) error message for this error
	InnerError *Error        // An underlying error object used to identify the root cause of this error.
	Details    []interface{} // Additional information related to this error message, such as parameters to the function that caused the error.
	TimeStamp  int64         // Unix Epoch timestamp of the date/time when this error was created
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

// Report sends this error to all configured reporters, to be reported via their various error reporting channels.
func (err *Error) Report() *Error {

	for _, reporter := range reporters {
		reporter.Report(err)
	}

	return err
}

// NotFound returns TRUE if the error `Code` is a 404 / Not Found error.
func (err *Error) NotFound() bool {
	return err.Code == CodeNotFound
}
