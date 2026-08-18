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

// TestAsError_TypedNilPointer confirms that a typed-nil *Error is not
// dereferenced, and converts into an empty (zero) Error instead.
func TestAsError_TypedNilPointer(t *testing.T) {

	var err *Error

	result := AsError(err)
	require.True(t, result.IsZero())
	require.Equal(t, 0, result.Code)
	require.Equal(t, "", result.Message)
}

// TestAsError_Pointer confirms that a populated *Error is dereferenced into its value.
func TestAsError_Pointer(t *testing.T) {

	inner := newError(codeNotFoundError, "Location", "Message")

	result := AsError(&inner)
	require.Equal(t, codeNotFoundError, result.Code)
	require.Equal(t, "Location", result.Location)
	require.Equal(t, "Message", result.Message)
}
