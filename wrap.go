package derp

import "time"

// Wrap encapsulates an existing derp.Error, and is guaranteed to return a "Not Nil" value.
// This function ALWAYS returns a non-nil error value.
func Wrap(inner error, location string, message string, details ...any) error {

	// double nil check to make nilaway happy.  NotNil also catches typed-nil
	// pointers, which would panic when Error() is called on them below.
	if (inner != nil) && NotNil(inner) {

		// If the inner error is not a derp.Error, then serialize it into the details.
		// derp.Errors are skipped because WrappedValue already preserves everything they carry.
		switch inner.(type) {
		case Error, *Error:
		default:
			// Clipped before appending: when the caller passes its own slice with `...`,
			// `details` aliases that slice, and a bare append would write into the
			// caller's spare capacity.
			details = append(details[:len(details):len(details)], inner.Error())
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
