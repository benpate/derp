package derp

import "time"

// AsError converts a generic error into a derp.Error type
func AsError(err error) Error {

	if err == nil {
		return Error{}
	}

	switch typed := err.(type) {

	case Error:
		return typed

	case *Error:
		return *typed

	default:
		return Error{
			Message:      "Wrapped non-derp error",
			WrappedValue: err,
			Location:     "derp.AsError",
			Code:         codeInternalError,
			TimeStamp:    time.Now().Unix(),
		}
	}
}
