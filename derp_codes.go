package derp

import (
	"reflect"
	"strings"
)

/******************************************
 * Error Code Functions
 * These determine if an error matches a specific error code.
 *****************************************/

// IsBadRequest returns TRUE if this is a 400 (Bad Request) error.
func IsBadRequest(err error) bool {
	return ErrorCode(err) == codeBadRequestError
}

// IsUnauthorized returns TRUE if this is a 401 (Unauthorized) error.
func IsUnauthorized(err error) bool {
	return ErrorCode(err) == codeUnauthorizedError
}

// IsForbidden returns TRUE if this is a 403 (Forbidden) error.
func IsForbidden(err error) bool {
	return ErrorCode(err) == codeForbiddenError
}

// IsNotFound returns TRUE if this is a 404 (Not Found) error.
func IsNotFound(err error) bool {
	if ErrorCode(err) == codeNotFoundError {
		return true
	}

	// Special case for "not found" string errors
	return strings.ToLower(err.Error()) == "not found"
}

// IsGone returns TRUE if this is a 410 (Gone) error.
func IsGone(err error) bool {
	return (ErrorCode(err) == codeGoneError)
}

// IsNotFoundOrGone returns TRUE if this is a 404 (Not Found) or 410 (Gone) error.
func IsNotFoundOrGone(err error) bool {

	switch ErrorCode(err) {
	case codeNotFoundError:
	case codeGoneError:
		return true
	}

	return (err.Error() == "not found")
}

// IsTeapot returns TRUE if this is a 418 (I'm a Teapot) error.
func IsTeapot(err error) bool {
	return ErrorCode(err) == codeTeapotError
}

// IsMisdirectedRequest returns TRUE if this is a 421 (Misdirected Request) error.
func IsMisdirectedRequest(err error) bool {
	return ErrorCode(err) == codeMisdirectedRequestError
}

// IsValidationError returns TRUE if this is a 422 (Validation) error.
func IsValidationError(err error) bool {
	return ErrorCode(err) == codeValidationError
}

// IsInternalServerError returns TRUE if this is a 500 (Internal Server Error) error.
func IsInternalServerError(err error) bool {
	return ErrorCode(err) == codeInternalError
}

// IsNotImplemented returns TRUE if this is a 501 (Not Implemented) error.
func IsNotImplemented(err error) bool {
	return ErrorCode(err) == codeNotImplementedError
}

/******************************************
 * Range Functions
 * These functions determine if an error is
 * within a certain range of HTTP status codes.
 *****************************************/

// IsInformational returns TRUE if the error `Code` is a 1xx / Informational error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#1xx_informational_response
func IsInformational(err error) bool {
	code := ErrorCode(err)
	return code >= 100 && code < 200
}

// IsSuccess returns TRUE if the error `Code` is a 2xx / Success error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#2xx_success
func IsSuccess(err error) bool {
	code := ErrorCode(err)
	return code >= 200 && code < 300
}

// IsRedirection returns TRUE if the error `Code` is a 3xx / Redirection error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#3xx_redirection
func IsRedirection(err error) bool {
	code := ErrorCode(err)
	return code >= 300 && code < 400
}

// IsClientError returns TRUE if the error `Code` is a 4xx / Client Error error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#4xx_client_errors
func IsClientError(err error) bool {
	code := ErrorCode(err)
	return code >= 400 && code < 500
}

// IsServerError returns TRUE if the error `Code` is a 5xx / Server Error error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#5xx_server_errors
func IsServerError(err error) bool {
	code := ErrorCode(err)
	return code >= 500 && code < 600
}

/******************************************
 * Other Utility Functions
 *****************************************/

// IsNil performs a robust nil check on an error interface
// Shout out to: https://medium.com/@mangatmodi/go-check-nil-interface-the-right-way-d142776edef1
func IsNil(err error) bool {

	if err == nil {
		return true
	}

	switch reflect.TypeOf(err).Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
		return reflect.ValueOf(err).IsNil()
	}

	return false
}

// NotNil returns TRUE if the error is NOT nil.
func NotNil(err error) bool {
	return !IsNil(err)
}
