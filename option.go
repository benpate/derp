package derp

// Option defines a function that modifies a derp.Error
type Option func(*Error)

func WithCode(code int) Option {
	return func(e *Error) {
		e.Code = code
	}
}

// WithBadRequest returns an option that sets the derp.Error code to 400 (Bad Request)
func WithBadRequest() Option {
	return WithCode(CodeBadRequestError)
}

// WithForbidden returns an option that sets the derp.Error code to 403 (Forbidden)
func WithForbidden() Option {
	return WithCode(CodeForbiddenError)
}

// WithNotFound returns an option that sets the derp.Error code to 404 (Not Found)
func WithNotFound() Option {
	return WithCode(CodeNotFoundError)
}

// WithInternalError returns an option that sets the derp.Error code to 500 (Internal Server Error)
func WithInternalError() Option {
	return WithCode(CodeInternalError)
}

// WithWrappedValue returns an option that sets the derp.Error wrapped value
func WithWrappedValue(inner error) Option {
	return func(e *Error) {
		e.WrappedValue = inner
	}
}

// WithLocation returns an option that sets the derp.Error location
func WithLocation(location string) Option {
	return func(e *Error) {
		e.Location = location
	}
}
