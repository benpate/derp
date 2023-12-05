package derp

type Option func(*SingleError)

func WithCode(code int) Option {
	return func(e *SingleError) {
		e.Code = code
	}
}

func WithNotFound() Option {
	return WithCode(CodeNotFoundError)
}

func WithBadRequest() Option {
	return WithCode(CodeBadRequestError)
}

func WithForbidden() Option {
	return WithCode(CodeForbiddenError)
}

func WithInternalError() Option {
	return WithCode(CodeInternalError)
}

func WithWrappedValue(inner error) Option {
	return func(e *SingleError) {
		e.WrappedValue = inner
	}
}

func WithLocation(location string) Option {
	return func(e *SingleError) {
		e.Location = location
	}
}
