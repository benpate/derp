package derp

// Error represents a Derp error
type Error struct {
	Code       int           // Numeric error code (such as an HTTP status code) to report to the client.
	TimeStamp  int64         // Date Time Stamp
	Location   string        // Function name (or other location description) of where the error occurred
	Message    string        // Primary (top-level) error message for this error
	Details    []interface{} // Additional information related to this error message
	InnerError *Error        // An underlying error object used to identify the root cause of this error.
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
