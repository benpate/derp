package derp

import (
	"encoding/json"
	"time"
)

/******************************************
 * New Derp
 ******************************************/

// BadRequest returns a (400) Bad Request error
// which indicates that the request is not properly formatted.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request
func BadRequest(location string, message string, details ...any) Error {
	return new(codeBadRequestError, location, message, details...)
}

// Unauthorized returns a (401) Unauthorized error
// which indicates that the request requires user authentication.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized
func Unauthorized(location string, message string, details ...any) Error {
	return new(codeUnauthorizedError, location, message, details...)
}

// Forbidden returns a (403) Forbidden error
// which indicates that the current user does not have permissions to access the requested resource.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden
func Forbidden(location string, message string, details ...any) Error {
	return new(codeForbiddenError, location, message, details...)
}

// NotFound returns a (404) Not Found error
// which indicates that the requested resource does not exist,
// such as when database query returns "not found"
// https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found
func NotFound(location string, message string, details ...any) Error {
	return new(codeNotFoundError, location, message, details...)
}

// Gone returns a (410) Gone error
// which indicates that the resource requested is no longer available and will not be available again.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-410-gone
func Gone(location string, message string, details ...any) Error {
	return new(codeGoneError, location, message, details...)
}

// Teapot returns a (418) I'm a Teapot error
// which indicates that the server is a teapot that cannot serve HTTP requests.
// https://www.rfc-editor.org/rfc/rfc7168.html#name-418-im-a-teapot
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#418
func Teapot(location string, message string, details ...any) Error {
	return new(codeTeapotError, location, message, details...)
}

// MisdirectedRequest returns a (421) Misdirected Request error.
// which indicates that the request was made to the wrong server; that server is not able to produce a response.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request
func MisdirectedRequest(location string, message string, details ...any) Error {
	return new(codeMisdirectedRequestError, location, message, details...)
}

// Validation returns a (422) Validation error
// which indicates that the request contains invalid data.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content
func Validation(message string, details ...any) Error {
	return new(codeValidationError, "", message, details...)
}

// Internal returns a (500) Internal Server Error
// which represents a generic error message, given when an unexpected condition was encountered and no more specific message is suitable.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error
func Internal(location string, message string, details ...any) Error {
	return new(codeInternalError, location, message, details...)
}

// NotImplemented returns a (501) Not Implemented error
// which indicates that the server does not support the functionality required to fulfill the request.
// https://www.rfc-editor.org/rfc/rfc9110.html#name-501-not-implemented
func NotImplemented(location string, details ...any) Error {
	return new(codeNotImplementedError, location, "Not Implemented", details...)
}

// Timeout returns a (524) Timeout error
// which indicates that the request took longer than an internal timeout threshold
// https://http.dev/524
func Timeout(location string, message string, details ...any) Error {
	return new(codeTimeout, location, message, details...)
}

/******************************************
 * Deprecated Derp (for backward compatibility)
 ******************************************/

// deprecated: use BadRequest() instead
func BadRequestError(location string, message string, details ...any) Error {
	return new(codeBadRequestError, location, message, details...)
}

// deprecated: use Unauthorized() instead
func UnauthorizedError(location string, message string, details ...any) Error {
	return new(codeUnauthorizedError, location, message, details...)
}

// deprecated: use Forbidden() instead
func ForbiddenError(location string, message string, details ...any) Error {
	return new(codeForbiddenError, location, message, details...)
}

// deprecated: use MisdirectedRequest() instead
func MisdirectedRequestError(location string, message string, details ...any) Error {
	return new(codeMisdirectedRequestError, location, message, details...)
}

// deprecated: use NotFound() instead
func NotFoundError(location string, message string, details ...any) Error {
	return new(codeNotFoundError, location, message, details...)
}

// deprecated: use Teapot() instead
func TeapotError(location string, message string, details ...any) Error {
	return new(codeTeapotError, location, message, details...)
}

// deprecated: use Timeout() instead
func TimeoutError(location string, message string, details ...any) Error {
	return new(codeTimeout, location, message, details...)
}

// deprecated: use Validation() instead
func ValidationError(message string, details ...any) Error {
	return new(codeValidationError, "", message, details...)
}

// deprecated: use Internal() instead
func InternalError(location string, message string, details ...any) Error {
	return new(codeInternalError, location, message, details...)
}

// deprecated: use NotImplemented() instead
func NotImplementedError(location string, details ...any) Error {
	return new(codeNotImplementedError, location, "Not Implemented", details...)
}

// new returns a new Error object
func new(code int, location string, message string, details ...any) Error {

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

	if IsNil(err) || err == nil {
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

	if err == nil || IsNil(err) {
		return 0
	}

	if getter, ok := err.(ErrorCodeGetter); ok {
		return getter.GetErrorCode()
	}

	return codeInternalError
}

func Location(err error) string {

	if IsNil(err) || err == nil {
		return ""
	}

	if getter, ok := err.(LocationGetter); ok {
		return getter.GetLocation()
	}

	return ""
}

func URL(err error) string {

	if IsNil(err) || err == nil {
		return ""
	}

	if getter, ok := err.(URLGetter); ok {
		return getter.GetURL()
	}

	return ""
}

func Details(err error) []any {

	if IsNil(err) || err == nil {
		return nil
	}

	if getter, ok := err.(DetailsGetter); ok {
		return getter.GetDetails()
	}

	return nil
}

func Serialize(err error) string {

	if NotNil(err) || err == nil {
		if bytes, err := json.Marshal(err); err == nil {
			return string(bytes)
		}
	}

	return ""
}

/******************************************
 * Root Values find the deeped properties available
 ******************************************/

// RootMessage returns the deepest message
// available within a chain of wrapped errors.
func RootMessage(err error) string {

	if IsNil(err) || err == nil {
		return ""
	}

	if unwrapper, isUnwrapper := err.(Unwrapper); isUnwrapper {
		wrappedError := unwrapper.Unwrap()
		wrappedMessage := Message(wrappedError)

		if wrappedMessage != "" {
			return wrappedMessage
		}
	}

	return Message(err)
}

// RootLocation returns the deepest location
// defined within a chain of wrapped errors.
func RootLocation(err error) string {

	if IsNil(err) || err == nil {
		return ""
	}

	if unwrapper, isUnwrapper := err.(Unwrapper); isUnwrapper {
		wrappedError := unwrapper.Unwrap()
		wrappedLocation := Location(wrappedError)

		if wrappedLocation != "" {
			return wrappedLocation
		}
	}

	return Location(err)
}
