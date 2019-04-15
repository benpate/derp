package derp

// Derp recommends, but does not require, using HTTP status codes as error messages.
// Several of the most useful messages are listed here as defaults.

const (
	// CodeNotFound represents a request for a resource that does not exist, such as a database query that returns "not found"
	CodeNotFound = 400

	// CodeForbidden means that the current user does not have the required permissions to access the requested resource.
	CodeForbidden = 403

	// CodeInternalServer represents a generic error message, given when an unexpected condition was encountered and no more specific message is suitable.
	CodeInternalError = 500
)
