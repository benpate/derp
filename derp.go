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
