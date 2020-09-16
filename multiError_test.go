package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiError_Append(t *testing.T) {

	{
		e := NewMultiError().Append(
			errors.New("first error here"),
			errors.New("second error here"),
			errors.New("third error here"),
		)

		require.Equal(t, 3, len(e.Errors))
		require.Equal(t, "first error here", e.Errors[0].Error())
		require.Equal(t, "second error here", e.Errors[1].Error())
		require.Equal(t, "third error here", e.Errors[2].Error())

		Report(e)
	}
}

func TestMultiError_AppendNested(t *testing.T) {

	{
		e := NewMultiError()
		e.Append(
			errors.New("first error here"),
			errors.New("second error here"),
			NewMultiError().Append(
				errors.New("first nested error here"),
				errors.New("second nested error here"),
			),
		)

		require.Equal(t, 4, len(e.Errors))
		require.Equal(t, 500, e.ErrorCode())
		require.Equal(t, 500, ErrorCode(e))
		require.Equal(t, "first error here", e.Errors[0].Error())
		require.Equal(t, "second error here", e.Errors[1].Error())
		require.Equal(t, "first nested error here", e.Errors[2].Error())
		require.Equal(t, "second nested error here", e.Errors[3].Error())

		Report(e)
	}
}

func TestMultiError_AppendNil(t *testing.T) {

	{
		e := NewMultiError().Append(nil, nil, nil)
		require.Zero(t, ErrorCode(e))
		require.Zero(t, len(e.Errors))
	}

	{
		e := &MultiError{}
		require.Zero(t, ErrorCode(e))
	}
}
