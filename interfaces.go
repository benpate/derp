package derp

import "time"

// DetailsGetter interface wraps the GetDetails method, which returns a list of details about the error
type DetailsGetter interface {
	// Details returns a list of details about the error.
	GetDetails() []any
}

// ErrorCodeGetter interface wraps the GetErrorCode method, which returns a numeric, application-specific code that references this error
type ErrorCodeGetter interface {

	// ErrorCode returns a numeric, application-specific code that references this error.
	// HTTP status codes are recommended, but not required
	GetErrorCode() int
}

// LocationGetter interface wraps the GetLocation method, which returns the location of the error
type LocationGetter interface {
	// Location returns the location of the error in the source code.
	GetLocation() string
}

// MessageGetter interface wraps the GetMessage method, which returns a human-friendly string representation of the error
type MessageGetter interface {

	// Message returns a human-friendly string representation of the error.
	GetMessage() string
}

// RetryAfterGetter interface wraps the GetRetryAfter method, which returns the number of seconds to wait before retrying the operation that caused this error
type RetryAfterGetter interface {
	// RetryAfter returns the number of seconds to wait before retrying the operation that caused this error.
	GetRetryAfter() time.Duration
}

// URLGetter interface wraps the GetURL method, which returns a URL to a web page with more information about this error
type URLGetter interface {
	// URL returns a URL to a web page with more information about this error.
	GetURL() string
}

// Unwrapper interface describes any error that can be "unwrapped".  It supports
// the Unwrap method added in Go 1.13+
type Unwrapper interface {

	// Unwrap returns the inner error bundled inside of an outer error.
	Unwrap() error
}
