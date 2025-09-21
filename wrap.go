package derp

import "time"

// Wrap encapsulates an existing derp.Error, and is guaranteed to return a "Not Nil" value.
// This function ALWAYS returns a non-nil error value.
func Wrap(inner error, location string, message string, details ...any) error {

	if inner != nil {
		// If the inner error is not of a known type, then serialize it into the details.
		switch inner.(type) {
		case Error:
		default:
			details = append(details, inner.Error())
		}
	}

	result := Error{
		WrappedValue: inner,
		Location:     location,
		Message:      message,
		Details:      make([]any, 0, len(details)),
		TimeStamp:    time.Now().Unix(),
		Code:         ErrorCode(inner),
	}

	for _, detail := range details {
		if option, ok := detail.(Option); ok {
			option(&result)
		} else {
			result.Details = append(result.Details, detail)
		}
	}

	return result
}

// WrapIF returns a wrapped error if the inner error is not nil.
// If the inner error is nil, then this function returns nil.
func WrapIF(inner error, location string, message string, details ...any) error {

	// If the inner error is nil, then the wrapped error is nil, too.
	if IsNil(inner) {
		return nil
	}

	return Wrap(inner, location, message, details...)
}
