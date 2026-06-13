package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAsError(t *testing.T) {

	// nil errors become the zero-value Error
	require.True(t, AsError(nil).IsZero())

	// a derp.Error value is returned unchanged
	{
		original := newError(404, "location", "message")
		require.Equal(t, original, AsError(original))
	}

	// a *derp.Error pointer is dereferenced
	{
		original := newError(404, "location", "message")
		require.Equal(t, original, AsError(&original))
	}

	// a standard error is wrapped as an Internal error
	{
		result := AsError(errors.New("standard error"))
		require.Equal(t, codeInternalError, result.Code)
		require.Equal(t, "derp.AsError", result.Location)
		require.Equal(t, "standard error", result.WrappedValue.Error())
	}
}
