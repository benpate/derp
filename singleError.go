package derp

// Error represents a runtime error.  It includes
type Error struct {
	Code         int    `json:"code"`                 // Numeric error code (such as an HTTP status code) to report to the client.
	Location     string `json:"location"`             // Function name (or other location description) of where the error occurred
	Message      string `json:"message"`              // Primary (top-level) error message for this error
	Details      []any  `json:"details,omitempty"`    // Additional information related to this error message, such as parameters to the function that caused the error.
	TimeStamp    int64  `json:"timestamp"`            // Unix Epoch timestamp of the date/time when this error was created
	WrappedValue error  `json:"innerError,omitempty"` // An underlying error object used to identify the root cause of this error.
}

// Error implements the Error interface, which allows derp.Error objects to be
// used anywhere a standard error is used.
func (err Error) Error() string {
	return err.Location + ": " + err.Message
}

// ErrorCode returns the error Code embedded in this Error.
func (err Error) GetErrorCode() int {
	return err.Code
}

// ErrorMessage returns the error Message embedded in this Error.
func (err Error) GetMessage() string {
	return err.Message
}

func (err Error) GetLocation() string {
	return err.Location
}

func (err Error) GetDetails() []any {
	return err.Details
}

func (err Error) GetTimeStamp() int64 {
	return err.TimeStamp
}

// Unwrap supports Go 1.13+ error unwrapping
func (err Error) Unwrap() error {
	return err.WrappedValue
}
