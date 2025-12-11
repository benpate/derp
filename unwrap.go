package derp

// RootCause digs into the error stack and returns the original error
// that caused the DERP.  This is an alias for the Unwrap() function.
func RootCause(err error) error {
	return Unwrap(err)
}

// Unwrap digs into the error stack and returns the original error
// that caused the DERP
func Unwrap(err error) error {

	// If this error can be "unwrapped" then dig deeper into the chain
	if unwrapper, ok := err.(Unwrapper); ok {

		// Try to unwrap the error.  If it is a not-Nil result, then keep digging
		if next := unwrapper.Unwrap(); NotNil(next) {
			return Unwrap(next)
		}
	}

	// Fall through means that there is nothing left to unwrap.  Return the current error
	return err
}

// UnwrapHTTPError unwraps the provided error, returning the first HTTPError found in the chain.
// If no HTTPError is found, this function returns nil.
func UnwrapHTTPError(err error) *HTTPError {

	if httpError, isHTTPError := err.(HTTPError); isHTTPError {
		return &httpError
	}

	// If this error can be "unwrapped" then dig deeper into the chain
	if unwrapper, ok := err.(Unwrapper); ok {

		// Try to unwrap the error.  If it is a not-Nil result, then keep digging
		if next := unwrapper.Unwrap(); NotNil(next) {
			return UnwrapHTTPError(next)
		}
	}

	// Fall through means that there is nothing left to unwrap.  Return the current error
	return nil

}
