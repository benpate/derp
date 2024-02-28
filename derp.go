package derp

import (
	"reflect"
	"time"
)

/******************************************
 * Constructor Functions
 ******************************************/

func NewBadRequestError(location string, message string, details ...any) Error {
	return New(CodeBadRequestError, location, message, details...)
}

func NewForbiddenError(location string, message string, details ...any) Error {
	return New(CodeForbiddenError, location, message, details...)
}

func NewInternalError(location string, message string, details ...any) Error {
	return New(CodeInternalError, location, message, details...)
}

func NewNotFoundError(location string, message string, details ...any) Error {
	return New(CodeNotFoundError, location, message, details...)
}

func NewUnauthorizedError(location string, message string, details ...any) Error {
	return New(CodeUnauthorizedError, location, message, details...)
}

func NewValidationError(message string, details ...any) Error {
	return New(CodeValidationError, "", message, details...)
}

// New returns a new Error object
func New(code int, location string, message string, details ...any) Error {

	result := Error{
		Location:  location,
		Code:      code,
		Message:   message,
		Details:   make([]any, 0, len(details)),
		TimeStamp: time.Now().Unix(),
	}

	for _, detail := range details {
		if option, ok := detail.(Option); ok {
			option(&result)
		} else {
			result.Details = append(result.Details, detail)
		}
	}

	return result
}

/******************************************
 * Data Accessor Functions
 ******************************************/

// Message retrieves the best-fit error message for any type of error
func Message(err error) string {

	if isNil(err) {
		return ""
	}

	if getter, ok := err.(MessageGetter); ok {
		return getter.GetMessage()
	}

	return err.Error()
}

// ErrorCode returns an error code for any error.  It tries to read the error code
// from objects matching the ErrorCodeGetter interface.  If the provided error does not
// match this interface, then it assigns a generic "Internal Server Error" code 500.
func ErrorCode(err error) int {

	if isNil(err) {
		return 0
	}

	if getter, ok := err.(ErrorCodeGetter); ok {
		return getter.GetErrorCode()
	}

	return CodeInternalError
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
	case Error:
	default:
		details = append(details, inner.Error())
	}

	result := Error{
		WrappedValue: inner,
		Location:     location,
		Message:      message,
		Details:      make([]any, 0, len(details)),
		TimeStamp:    time.Now().Unix(),
		Code:         ErrorCode(inner),
	}

	for _, detail := range details {
		if option, ok := detail.(Option); ok {
			option(&result)
		} else {
			result.Details = append(result.Details, detail)
		}
	}

	return result
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
