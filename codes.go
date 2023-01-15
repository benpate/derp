package derp

// Derp recommends, but does not require, using HTTP status codes as error messages.
// Several of the most useful messages are listed here as defaults.

const (

	// CodeBadRequestError indicates that the request is not properly formatted.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request
	CodeBadRequestError = 400

	// CodeUnauthorizedError indicates that the request requires user authentication.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized
	CodeUnauthorizedError = 401

	// CodeForbiddenError means that the current user does not have the required permissions to access the requested resource.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden
	CodeForbiddenError = 403

	// CodeNotFoundError represents a request for a resource that does not exist, such as a database query that returns "not found"
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found
	CodeNotFoundError = 404

	// CodeInternalError represents a generic error message, given when an unexpected condition was encountered and no more specific message is suitable.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error
	CodeInternalError = 500

	// CodeValidationError retpresents a request that contains invalid data.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content
	CodeValidationError = 422
)
