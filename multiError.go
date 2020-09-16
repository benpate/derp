package derp

import "strings"

// MultiError represents a runtime error.  It includes
type MultiError struct {
	Errors []error
}

// Error implements the Error interface, which allows derp.Error objects to be
// used anywhere a standard error is used.
func (err *MultiError) Error() string {

	b := strings.Builder{}

	for _, e := range err.Errors {
		b.WriteString(e.Error() + "\n")
	}

	return b.String()
}

// ErrorCode returns the error Code embedded in this Error.  This is useful for matching
// interfaces in other package.
func (err *MultiError) ErrorCode() int {

	if len(err.Errors) == 0 {
		return 0
	}

	return ErrorCode(err.Errors[0])
}
