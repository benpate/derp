package derp

import (
	"strings"
)

// MultiError represents a runtime error.  It includes
type MultiError []error

func (m *MultiError) Append(err error) {
	if !isNil(err) {
		*m = append(*m, err)
	}
}

func (m MultiError) Length() int {
	if isNil(m) {
		return 0
	}
	return len(m)
}

func (m MultiError) IsEmpty() bool {
	return m.Length() == 0
}

// Message retrieves the error message from the first message in the slice that is a messageGetter
func (m MultiError) Message() string {

	for _, err := range m {

		if message := Message(err); message != "" {
			return message
		}
	}

	return ""
}

func (m MultiError) AddPrefixes(prefix string) {

	for index, err := range m {
		switch typed := err.(type) {
		case SingleError:
			typed.SetMessage(prefix + typed.Message)
			m[index] = typed
		case MessageSetter:
			typed.SetMessage(prefix + Message(err))
		}
	}
}

// Error implements the Error interface, which allows derp.Error objects to be
// used anywhere a standard error is used.
func (m MultiError) Error() string {

	b := strings.Builder{}

	for _, err := range m {
		b.WriteString(err.Error() + "\n")
	}

	return b.String()
}

// ErrorCode returns the error Code embedded in this Error.  This is useful for matching
// interfaces in other package.
func (m MultiError) ErrorCode() int {

	if m.IsEmpty() {
		return 0
	}

	for _, err := range m {

		if errorWithCode, ok := err.(ErrorCodeGetter); ok {
			return errorWithCode.ErrorCode()
		}
	}

	return 500
}
