package derp

import "time"

// Error represents a runtime error, including a numeric code, the location
// where it occurred, a human-readable message, and an optional wrapped error.
type Error struct {
	Code         int    `json:"code"`                 // Numeric error code (such as an HTTP status code) to report to the client.
	Location     string `json:"location"`             // Function name (or other location description) of where the error occurred
	Message      string `json:"message"`              // Primary (top-level) error message for this error
	URL          string `json:"url,omitempty"`        // URL to a web page with more information about this error
	Details      []any  `json:"details,omitempty"`    // Additional information related to this error message, such as parameters to the function that caused the error.
	TimeStamp    int64  `json:"timestamp"`            // Unix Epoch timestamp of the date/time when this error was created
	WrappedValue error  `json:"innerError,omitempty"` // An underlying error object used to identify the root cause of this error.
}

// IsZero returns true if this Error is empty / uninitialized
func (err Error) IsZero() bool {
	if err.Code != 0 {
		return false
	}

	if err.Location != "" {
		return false
	}

	if err.Message != "" {
		return false
	}

	if err.WrappedValue != nil {
		return false
	}

	if err.URL != "" {
		return false
	}

	if len(err.Details) > 0 {
		return false
	}

	return true
}

// Error implements the Error interface, which allows derp.Error objects to be
// used anywhere a standard error is used.
func (err Error) Error() string {
	return err.Location + ": " + err.Message
}

// GetErrorCode returns the error Code embedded in this Error.
func (err Error) GetErrorCode() int {
	return err.Code
}

// GetLocation returns the error Location embedded in this Error.
func (err Error) GetLocation() string {
	return err.Location
}

// GetMessage returns the error Message embedded in this Error.
func (err Error) GetMessage() string {
	return err.Message
}

// GetRetryAfter returns the retry-after duration carried by this Error's WrappedValue,
// or zero when nothing below it carries one.
func (err Error) GetRetryAfter() time.Duration {

	// RULE: This MUST stay a delegation.  RetryAfter's errors.As matches this Error first,
	// so re-entering it one level down is what searches the rest of the chain.  Reading a
	// field here instead would stop the search at the top.
	return RetryAfter(err.WrappedValue)
}

// GetURL returns the help URL embedded in this Error.
func (err Error) GetURL() string {
	return err.URL
}

// GetDetails returns the error Details embedded in this Error.
func (err Error) GetDetails() []any {
	return err.Details
}

// GetTimeStamp returns the error TimeStamp embedded in this Error.
func (err Error) GetTimeStamp() int64 {
	return err.TimeStamp
}

// Unwrap supports Go 1.13+ error unwrapping
func (err Error) Unwrap() error {
	return err.WrappedValue
}
