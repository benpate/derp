package derp

import (
	"time"
)

/******************************************
 * Constructor Functions
 ******************************************/

// NewBadRequestError returns a (400) Bad Request error
func NewBadRequestError(location string, message string, details ...any) Error {
	return New(CodeBadRequestError, location, message, details...)
}

// NewUnauthorizedError returns a (401) Unauthorized error
func NewUnauthorizedError(location string, message string, details ...any) Error {
	return New(CodeUnauthorizedError, location, message, details...)
}

// NewForbiddenError returns a (403) Forbidden error
func NewForbiddenError(location string, message string, details ...any) Error {
	return New(CodeForbiddenError, location, message, details...)
}

// NewNotFoundError returns a (404) Not Found error
func NewNotFoundError(location string, message string, details ...any) Error {
	return New(CodeNotFoundError, location, message, details...)
}

// NewTeapotError returns a (418) I'm a Teapot error
func NewTeapotError(location string, message string, details ...any) Error {
	return New(CodeTeapotError, location, message, details...)
}

// NewMisdirectedRequestError returns a (421) Misdirected Request error.
func NewMisdirectedRequestError(location string, message string, details ...any) Error {
	return New(CodeMisdirectedRequestError, location, message, details...)
}

// NewValidationError returns a (422) Validation error
func NewValidationError(message string, details ...any) Error {
	return New(CodeValidationError, "", message, details...)
}

// NewInternalError returns a (500) Internal Server Error
func NewInternalError(location string, message string, details ...any) Error {
	return New(CodeInternalError, location, message, details...)
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
