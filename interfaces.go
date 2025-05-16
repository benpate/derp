package derp

// ErrorCodeGetter interface describes any error that can also "get" an error code value
type ErrorCodeGetter interface {

	// ErrorCode returns a numeric, application-specific code that references this error.
	// HTTP status codes are recommended, but not required
	GetErrorCode() int
}

// MessageGetter interface describes any error that can also report a "Message"
type MessageGetter interface {

	// Message returns a human-friendly string representation of the error.
	GetMessage() string
}

type LocationGetter interface {
	// Location returns the location of the error in the source code.
	GetLocation() string
}

type DetailsGetter interface {
	// Details returns a list of details about the error.
	GetDetails() []any
}

// Unwrapper interface describes any error that can be "unwrapped".  It supports
// the Unwrap method added in Go 1.13+
type Unwrapper interface {

	// Unwrap returns the inner error bundled inside of an outer error.
	Unwrap() error
}
