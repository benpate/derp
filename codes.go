package derp

// Derp recommends, but does not require, using HTTP status codes as error messages.
// The values used in derp functions are enumerated here.

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

	// CodeTeapotError indicates that the server is a teapot and cannot serve HTTP requests.
	// https://www.rfc-editor.org/rfc/rfc7168.html#name-418-im-a-teapot
	// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#418
	CodeTeapotError = 418

	// CodeMisdirectedRequestError indicates that the request was directed to a server that is not able to produce a response.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request
	CodeMisdirectedRequestError = 421

	// CodeValidationError represents a request that contains invalid data.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content
	CodeValidationError = 422

	// CodeInternalError represents a generic error message, given when an unexpected condition was encountered and no more specific message is suitable.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error
	CodeInternalError = 500
)
