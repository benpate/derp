package derp

import (
	"reflect"
	"time"
)

// New returns a new Error object
func New(code int, location string, message string, details ...interface{}) *SingleError {

	return &SingleError{
		Location:  location,
		Code:      code,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().Truncate(1 * time.Second),
	}
}

// Wrap encapsulates an existing derp.Error
func Wrap(inner error, location string, message string, details ...interface{}) *SingleError {

	// If the inner error is nil, then the wrapped error is nil, too.
	if isNil(inner) {
		return nil
	}

	// If the inner error is not of a known type, then serialize it into the details.
	switch inner.(type) {
	case *SingleError:
	case *MultiError:
	case *ValidationError:
	default:
		details = append(details, inner.Error())
	}

	return &SingleError{
		InnerError: inner,
		Location:   location,
		Message:    message,
		Details:    details,
		TimeStamp:  time.Now().Truncate(1 * time.Second),
		Code:       ErrorCode(inner),
	}
}

// Append takes one or more errors and combines them into
// this MultiError.  If one of the arguments is itself a
// MultiError, then its slice of errors is flattened into
// this one, so that the final result is a single, one-dimensional
// slice of errors.
func Append(errs ...error) *MultiError {

	// TODO: Handle a list of all nils.

	result := &MultiError{
		Errors: make([]error, 0),
	}

	for _, e := range errs {

		// If this entry is nil, then skip it
		if isNil(e) {
			continue
		}

		// If this entry is a MultiError of its own, then flatten it into the new error
		if multi, ok := e.(*MultiError); ok {
			result.Errors = append(result.Errors, multi.Errors...)
			continue
		}

		// Otherwise, append the error into the resut
		result.Errors = append(result.Errors, e)
	}

	// If there is at least one error in the list, then this is a real MultiError
	if len(result.Errors) > 0 {
		return result
	}

	// Otherwise, this is "empty" so return nil.
	return nil
}

// NotFound returns TRUE if the error `Code` is a 404 / Not Found error.
func NotFound(err error) bool {

	if coder, ok := err.(ErrorCodeGetter); ok {
		return coder.ErrorCode() == CodeNotFoundError
	}

	return err.Error() == "not found"
}

func isNil(i error) bool {

	// Shout out to: https://medium.com/@mangatmodi/go-check-nil-interface-the-right-way-d142776edef1
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Array, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
