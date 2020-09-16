package derp

// Derp recommends, but does not require, using HTTP status codes as error messages.
// Several of the most useful messages are listed here as defaults.

const (

	// CodeBadRequestError indicates that the request is not properly formatted.
	CodeBadRequestError = 400

	// CodeForbiddenError means that the current user does not have the required permissions to access the requested resource.
	CodeForbiddenError = 403

	// CodeNotFoundError represents a request for a resource that does not exist, such as a database query that returns "not found"
	CodeNotFoundError = 404

	// CodeInternalError represents a generic error message, given when an unexpected condition was encountered and no more specific message is suitable.
	CodeInternalError = 500
)

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
