package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestError(t *testing.T) {
	e := newError(codeInternalError, "Location", "Message")
	require.Equal(t, "Location: Message", e.Error())
	require.Equal(t, codeInternalError, ErrorCode(e))

	WithNotFound()(&e)
	require.Equal(t, 404, e.GetErrorCode())
	require.Equal(t, 404, ErrorCode(e))
	require.True(t, IsNotFound(e))
}

func TestError_WrapSingle(t *testing.T) {

	inner := newError(101, "A", "B", "C")
	outer := Wrap(inner, "C", "D").(Error)

	require.Equal(t, outer.Code, 101)
	require.Equal(t, outer.Location, "C")
	require.Equal(t, outer.Message, "D")

	innerAgain := outer.Unwrap().(Error)
	require.Equal(t, innerAgain.Code, 101)
	require.Equal(t, innerAgain.Location, "A")
	require.Equal(t, innerAgain.Message, "B")
}

func TestError_WrapGeneric(t *testing.T) {

	inner := errors.New("omg it works")
	outer := Wrap(inner, "C", "D").(Error)

	require.Equal(t, outer.Code, 500)
	require.Equal(t, outer.Location, "C")
	require.Equal(t, outer.Message, "D")

	innerAgain := outer.Unwrap()
	require.Equal(t, innerAgain.Error(), "omg it works")
}

func TestErrorCodeSetter(t *testing.T) {
	err := Internal("test", "test", WithNotFound())
	require.Equal(t, 404, ErrorCode(err))
}

func TestIsZero(t *testing.T) {

	// the zero value is zero
	require.True(t, Error{}.IsZero())

	// any populated field makes it non-zero
	require.False(t, Error{Code: 500}.IsZero())
	require.False(t, Error{Location: "location"}.IsZero())
	require.False(t, Error{Message: "message"}.IsZero())
	require.False(t, Error{URL: "https://example.com"}.IsZero())
	require.False(t, Error{Details: []any{"detail"}}.IsZero())
	require.False(t, Error{WrappedValue: Error{}}.IsZero())
}

func TestErrorAccessors(t *testing.T) {

	err := Error{
		Code:      404,
		Location:  "the location",
		Message:   "the message",
		URL:       "https://example.com/help",
		Details:   []any{"a", "b"},
		TimeStamp: 1234567890,
	}

	require.Equal(t, 404, err.GetErrorCode())
	require.Equal(t, "the location", err.GetLocation())
	require.Equal(t, "the message", err.GetMessage())
	require.Equal(t, "https://example.com/help", err.GetURL())
	require.Equal(t, []any{"a", "b"}, err.GetDetails())
	require.Equal(t, int64(1234567890), err.GetTimeStamp())
}
