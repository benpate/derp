package derp

// ValidationError represents an input validation error, and includes fields necessary to
// report problems back to the end user.
type ValidationError struct {
	Path    string `json:"path"`    // Identifies the PATH (or variable name) that has invalid input
	Message string `json:"message"` // Human-readable message that explains the problem with the input value.
}

// Invalid returns a fully populated ValidationError to the caller
func Invalid(path string, message string) error {
	return &ValidationError{
		Path:    path,
		Message: message,
	}
}

// Error returns a string representation of this ValidationError, and implements
// the builtin errors.error interface.
func (v *ValidationError) Error() string {
	return v.Message
}

// ErrorCode returns CodeValidationError for this ValidationError
// It implements the ErrorCodeGetter interface.
func (v *ValidationError) ErrorCode() int {
	return CodeValidationError
}
