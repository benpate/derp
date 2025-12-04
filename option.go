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
	return WithCode(codeBadRequestError)
}

// WithForbidden returns an option that sets the derp.Error code to 403 (Forbidden)
func WithForbidden() Option {
	return WithCode(codeForbiddenError)
}

// WithNotFound returns an option that sets the derp.Error code to 404 (Not Found)
func WithNotFound() Option {
	return WithCode(codeNotFoundError)
}

// WithInternalError returns an option that sets the derp.Error code to 500 (Internal Server Error)
func WithInternalError() Option {
	return WithCode(codeInternalError)
}

// WithUnauthorized returns an option that sets the derp.Error code to 401 (Unauthorized)
func WithUnauthorized() Option {
	return WithCode(codeUnauthorizedError)
}

// WithWrappedValue returns an option that sets the derp.Error wrapped value
func WithWrappedValue(inner error) Option {
	return func(e *Error) {
		e.WrappedValue = inner
	}
}

// WithMessage returns an option that sets the derp.Error message
func WithMessage(message string) Option {
	return func(e *Error) {
		e.Message = message
	}
}

// WithLocation returns an option that sets the derp.Error location
func WithLocation(location string) Option {
	return func(e *Error) {
		e.Location = location
	}
}
