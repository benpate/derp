package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSingleError(t *testing.T) {
	e := New(CodeInternalError, "Location", "Message")
	require.Equal(t, "Location: Message", e.Error())
	require.Equal(t, CodeInternalError, ErrorCode(e))

	e.SetErrorCode(404)
	require.Equal(t, 404, e.ErrorCode())
	require.Equal(t, 404, ErrorCode(e))
	require.True(t, NotFound(e))
}

func TestSingleError_WrapSingle(t *testing.T) {

	inner := New(101, "A", "B", "C")
	outer := Wrap(inner, "C", "D")

	require.Equal(t, outer.Code, 101)
	require.Equal(t, outer.Location, "C")
	require.Equal(t, outer.Message, "D")

	innerAgain := outer.Unwrap().(*SingleError)
	require.Equal(t, innerAgain.Code, 101)
	require.Equal(t, innerAgain.Location, "A")
	require.Equal(t, innerAgain.Message, "B")
}

func TestSingleError_WrapMultiple(t *testing.T) {

	inner := NewMultiError().Append(
		errors.New("first error"),
		errors.New("second error"),
		errors.New("third error"),
	)

	outer := Wrap(inner, "C", "D")

	require.Equal(t, outer.Code, 500)
	require.Equal(t, outer.Location, "C")
	require.Equal(t, outer.Message, "D")

	innerAgain := outer.Unwrap().(*MultiError)
	require.Equal(t, innerAgain.Errors[0].Error(), "first error")
	require.Equal(t, innerAgain.Errors[1].Error(), "second error")
	require.Equal(t, innerAgain.Errors[2].Error(), "third error")
}

func TestSingleError_WrapGeneric(t *testing.T) {

	inner := errors.New("omg it works")
	outer := Wrap(inner, "C", "D")

	require.Equal(t, outer.Code, 500)
	require.Equal(t, outer.Location, "C")
	require.Equal(t, outer.Message, "D")

	innerAgain := outer.Unwrap()
	require.Equal(t, innerAgain.Error(), "omg it works")
}
