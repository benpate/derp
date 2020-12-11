package derp

import (
	"reflect"
	"time"
)

// New returns a new Error object
func New(code int, location string, message string, details ...interface{}) SingleError {

	return SingleError{
		Location:  location,
		Code:      code,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().Truncate(1 * time.Second),
	}
}

// Wrap encapsulates an existing derp.Error
func Wrap(inner error, location string, message string, details ...interface{}) error {

	// If the inner error is nil, then the wrapped error is nil, too.
	if isNil(inner) {
		return nil
	}

	// If the inner error is not of a known type, then serialize it into the details.
	switch inner.(type) {
	case SingleError:
	case MultiError:
	default:
		details = append(details, inner.Error())
	}

	return SingleError{
		InnerError: inner,
		Location:   location,
		Message:    message,
		Details:    details,
		TimeStamp:  time.Now().Truncate(1 * time.Second),
		Code:       ErrorCode(inner),
	}
}

// Message retrieves the best-fit error message for any type of error
func Message(err error) string {

	switch d := err.(type) {
	case *SingleError:
		return d.Message
	case *MultiError:
		if len(*d) > 0 {
			return Message((*d)[0])
		}
	}

	return err.Error()
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
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
