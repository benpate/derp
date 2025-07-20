package derp

// Derp recommends, but does not require, using HTTP status codes as error messages.
// The values used in derp functions are enumerated here.

const (

	// codeBadRequestError indicates that the request is not properly formatted.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-400-bad-request
	codeBadRequestError = 400

	// codeUnauthorizedError indicates that the request requires user authentication.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-401-unauthorized
	codeUnauthorizedError = 401

	// codeForbiddenError means that the current user does not have the required permissions to access the requested resource.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-403-forbidden
	codeForbiddenError = 403

	// codeNotFoundError represents a request for a resource that does not exist, such as a database query that returns "not found"
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-404-not-found
	codeNotFoundError = 404

	// codeTeapotError indicates that the server is a teapot and cannot serve HTTP requests.
	// https://www.rfc-editor.org/rfc/rfc7168.html#name-418-im-a-teapot
	// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#418
	codeTeapotError = 418

	// codeMisdirectedRequestError indicates that the request was directed to a server that is not able to produce a response.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-421-misdirected-request
	codeMisdirectedRequestError = 421

	// codeValidationError represents a request that contains invalid data.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-422-unprocessable-content
	codeValidationError = 422

	// codeInternalError represents a generic error message, given when an unexpected condition was encountered and no more specific message is suitable.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-500-internal-server-error
	codeInternalError = 500

	// codeNotImplementedError indicates that the server does not support the functionality required to fulfill the request.
	// https://www.rfc-editor.org/rfc/rfc9110.html#name-501-not-implemented
	codeNotImplementedError = 501

	// codeTimout is an unofficial server error, used by derp to indicate an internal timeout.
	// https://http.dev/524
	codeTimeout = 524
)
