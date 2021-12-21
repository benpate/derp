package derp

// Append collects errors into a multi-error, while still working nicely with nil errors.
func Append(original error, new error) error {

	if original == nil {
		return new
	}

	if new == nil {
		return original
	}

	switch o := original.(type) {
	case MultiError:
		return append(o, new)
	default:
		return MultiError{original, new}
	}
}
