package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// Tests for Wrap and WrapIF in wrap.go

func TestWrapGenericError(t *testing.T) {

	generic := errors.New("oof. that was bad")
	err := Wrap(generic, "TestEmptyWrappedValue", "Don't Do This").(Error)

	assert.Equal(t, 500, err.Code)
	assert.NotNil(t, err.WrappedValue)
	assert.Equal(t, "TestEmptyWrappedValue", err.Location)
	assert.Equal(t, "Don't Do This", err.Message)
	// assert.Equal(t, len(err.Details), 1)

	unwrapped := err.Unwrap()
	assert.Equal(t, "oof. that was bad", unwrapped.Error())
	Report(err)
}

func TestWrap_EmptyValue(t *testing.T) {

	{
		err := Wrap(nil, "TestEmptyWrappedValue", "This will still return an error")
		assert.Error(t, err)
	}

	{
		var innerError error
		outer := Wrap(innerError, "Should Still Return an error value", "really")
		assert.Error(t, outer)
	}
}

func TestWrapIF_EmptyValue(t *testing.T) {

	{
		err := WrapIF(nil, "TestEmptyWrappedValue", "This should return nil")
		assert.Nil(t, err)
	}

	{
		var innerError error
		outer := WrapIF(innerError, "Should Still Be Empty", "Really")
		assert.Nil(t, outer)
	}
}

func TestWrapIF_NotNil(t *testing.T) {
	err := WrapIF(errors.New("inner"), "location", "message")
	require.Error(t, err)
	require.Equal(t, codeInternalError, ErrorCode(err))
}

/******************************************
 * Fuzz Tests
 ******************************************/

// TestWrap_TypedNilPointer confirms that wrapping a typed-nil *Error does not
// call Error() on the nil pointer.
func TestWrap_TypedNilPointer(t *testing.T) {

	var inner *Error

	// Wrap always returns a non-nil error, even around a typed-nil inner value
	wrapped := Wrap(inner, "Location", "Message")
	require.NotNil(t, wrapped)
	require.Equal(t, "Location: Message", wrapped.Error())

	// The typed-nil inner value is not serialized into the Details
	require.Empty(t, Details(wrapped))

	// WrapIF treats a typed-nil inner value as nil, and returns nil
	require.Nil(t, WrapIF(inner, "Location", "Message"))
}

// FuzzWrapChain confirms that arbitrarily deep chains of wrapped errors can be
// unwrapped without panicking, and that the root values come from the innermost error.
func FuzzWrapChain(f *testing.F) {

	f.Add(uint8(0), "root", "outer")
	f.Add(uint8(1), "", "")
	f.Add(uint8(10), "root message", "wrapper")
	f.Add(uint8(255), "deep", "d")

	f.Fuzz(func(t *testing.T, depth uint8, rootMessage string, wrapMessage string) {

		// Build a chain of `depth` wrappers around a single root error
		root := newError(codeNotFoundError, "RootLocation", rootMessage)

		var err error = root
		for index := 0; index < int(depth); index++ {
			err = Wrap(err, "WrapLocation", wrapMessage)
		}

		// Wrap NEVER returns nil
		require.NotNil(t, err)

		// The code is inherited all the way up the chain
		require.Equal(t, codeNotFoundError, ErrorCode(err))

		// Unwrapping always terminates at the root error
		require.Equal(t, root.Error(), Unwrap(err).Error())
		require.Equal(t, RootCause(err), Unwrap(err))

		// The root location is always the innermost location
		require.Equal(t, "RootLocation", RootLocation(err))

		// The root message is the innermost NON-EMPTY message, so an empty root
		// message falls back to the message of the shallowest wrapper that has one.
		if rootMessage != "" {
			require.Equal(t, rootMessage, RootMessage(err))
		}
	})
}

// TestWrap_DoesNotAliasCallerSlice confirms that Wrap never writes into the spare
// capacity of a details slice that the caller passed with `...`.
func TestWrap_DoesNotAliasCallerSlice(t *testing.T) {

	// A slice with spare capacity, which the caller still owns
	backing := make([]any, 1, 4)
	backing[0] = "first"

	_ = Wrap(errors.New("inner-message"), "location", "message", backing...)

	// Re-slicing into the caller's spare capacity must NOT reveal derp's data
	require.Nil(t, backing[:2][1])
}

// TestWrap_PointerToError confirms that a *derp.Error is recognized as a derp error,
// and is not re-serialized into the Details of the wrapper.
func TestWrap_PointerToError(t *testing.T) {

	inner := NotFound("inner-location", "inner-message")
	wrapped := Wrap(&inner, "outer-location", "outer-message")

	require.Empty(t, AsError(wrapped).Details)
	require.Equal(t, codeNotFoundError, ErrorCode(wrapped))
}
