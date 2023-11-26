package derp

import (
	"reflect"
	"time"
)

/******************************************
 * Constructor Functions
 ******************************************/

func NewBadRequestError(location string, message string, details ...any) SingleError {
	return New(CodeBadRequestError, location, message, details...)
}

func NewForbiddenError(location string, message string, details ...any) SingleError {
	return New(CodeForbiddenError, location, message, details...)
}

func NewInternalError(location string, message string, details ...any) SingleError {
	return New(CodeInternalError, location, message, details...)
}

func NewNotFoundError(location string, message string, details ...any) SingleError {
	return New(CodeNotFoundError, location, message, details...)
}

func NewUnauthorizedError(location string, message string, details ...any) SingleError {
	return New(CodeUnauthorizedError, location, message, details...)
}

func NewValidationError(message string, details ...any) SingleError {
	return New(CodeValidationError, "", message, details...)
}

// New returns a new Error object
func New(code int, location string, message string, details ...any) SingleError {

	return SingleError{
		Location:  location,
		Code:      code,
		Message:   message,
		Details:   details,
		TimeStamp: time.Now().Unix(),
	}
}

/******************************************
 * Data Accessor Functions
 ******************************************/

// Message retrieves the best-fit error message for any type of error
func Message(err error) string {

	if isNil(err) {
		return ""
	}

	switch e := err.(type) {
	case SingleError:
		return e.Error()

	case MessageGetter:
		return e.Message()
	}

	return err.Error()
}

// SetMessage sets the error message for any errors that allow it.
func SetMessage(err error, message string) {

	if isNil(err) {
		return
	}

	switch e := err.(type) {
	case SingleError:
		e.SetMessage(message)
	case MessageSetter:
		e.SetMessage(message)
	}
}

// ErrorCode returns an error code for any error.  It tries to read the error code
// from objects matching the ErrorCodeGetter interface.  If the provided error does not
// match this interface, then it assigns a generic "Internal Server Error" code 500.
func ErrorCode(err error) int {

	if isNil(err) {
		return 0
	}

	if getter, ok := err.(ErrorCodeGetter); ok {
		return getter.ErrorCode()
	}

	return CodeInternalError
}

// SetErrorCode tries to set an error code for the provided error.  If the error matches the
// ErrorCodeSetter interface, then the code is set directly in the error.  Otherwise,
// it has no effect.
func SetErrorCode(err error, code int) {

	if isNil(err) {
		return
	}

	if setter, ok := err.(ErrorCodeSetter); ok {
		setter.SetErrorCode(code)
	}
}

// NotFound returns TRUE if the error `Code` is a 404 / Not Found error.
func NotFound(err error) bool {

	if isNil(err) {
		return false
	}

	if coder, ok := err.(ErrorCodeGetter); ok {
		return coder.ErrorCode() == CodeNotFoundError
	}

	return err.Error() == "not found"
}

// NilOrNotFound returns TRUE if the error is nil or a 404 / Not Found error.
// All other errors return FALSE
func NilOrNotFound(err error) bool {

	if isNil(err) {
		return true
	}

	if NotFound(err) {
		return true
	}

	return false
}

/******************************************
 * Other Manipulations
 ******************************************/

// Wrap encapsulates an existing derp.Error
func Wrap(inner error, location string, message string, details ...any) error {

	// If the inner error is nil, then the wrapped error is nil, too.
	if isNil(inner) {
		return nil
	}

	// If the inner error is not of a known type, then serialize it into the details.
	switch inner.(type) {
	case SingleError:
	default:
		details = append(details, inner.Error())
	}

	return SingleError{
		InnerError: inner,
		Location:   location,
		Message:    message,
		Details:    details,
		TimeStamp:  time.Now().Unix(),
		Code:       ErrorCode(inner),
	}
}

// ReportAndReturn reports an error to the logger
// and also returns it to the caller.
func ReportAndReturn(err error) error {
	Report(err)
	return err
}

/******************************************
 * Other Helpers
 ******************************************/

// isNil performs a robust nil check on an error interface
// Shout out to: https://medium.com/@mangatmodi/go-check-nil-interface-the-right-way-d142776edef1
func isNil(i error) bool {

	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
