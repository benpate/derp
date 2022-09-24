package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDerp tests basic derp functions (separate from features of a specific reporter)
func TestDerp(t *testing.T) {

	// Create an inner error
	innerError := New(CodeNotFoundError, "InnerError", "Not Found", "detail1", "detail2", "detail3")

	assert.Equal(t, innerError.Location, "InnerError")
	assert.Equal(t, innerError.Message, "Not Found")
	assert.Equal(t, innerError.Code, 404)
	assert.Equal(t, innerError.Details[0], "detail1")
	assert.Equal(t, innerError.Details[1], "detail2")
	assert.Equal(t, innerError.Details[2], "detail3")
	assert.Equal(t, NotFound(innerError), true)

	// Create an outer error
	outerError := Wrap(innerError, "OuterError", "Inherited", "other details here").(SingleError)

	assert.Equal(t, outerError.Location, "OuterError")
	assert.Equal(t, outerError.Message, "Inherited")
	assert.Equal(t, outerError.Code, 404) // This is still 404 because we've let the inner error code bubble up
	assert.NotNil(t, outerError.InnerError)
	assert.Equal(t, outerError.Details[0], "other details here")
	assert.Equal(t, NotFound(outerError), true)

	// Test the RootCause() function
	assert.Equal(t, "InnerError", RootCause(outerError).(SingleError).Location)
}

func TestConvenienceFns(t *testing.T) {

	badRequest := NewBadRequestError("location", "description")
	require.Equal(t, CodeBadRequestError, ErrorCode(badRequest))

	forbidden := NewForbiddenError("location", "description")
	require.Equal(t, CodeForbiddenError, ErrorCode(forbidden))

	internal := NewInternalError("location", "description")
	require.Equal(t, CodeInternalError, ErrorCode(internal))

	notFound := NewNotFoundError("location", "description")
	require.Equal(t, CodeNotFoundError, ErrorCode(notFound))

	unauthorized := NewUnauthorizedError("location", "description")
	require.Equal(t, CodeUnauthorizedError, ErrorCode(unauthorized))
}
func TestErrorInterface(t *testing.T) {

	// Create an error
	innerError := New(CodeNotFoundError, "Location Name", "Error Description", "details")

	// Verify that the error interface is outputting what we expect.
	assert.Equal(t, innerError.Error(), "Location Name: Error Description")
}

func TestStandardError(t *testing.T) {

	// Testing how derp handles an error from the standard library
	err := errors.New("This is a standard error")

	// Wrap it the stdlib error in a derp.  This means: 1) assigning an error code, and 2) making the original error message a property of the derp.Error.
	outer := New(CodeInternalError, "TestStandardError", "Encapsulating Error", err.Error())

	assert.Equal(t, "TestStandardError", outer.Location)
	assert.Equal(t, "Encapsulating Error", outer.Message)
	assert.Equal(t, CodeInternalError, outer.Code)
	assert.Equal(t, 1, len(outer.Details))
	assert.Nil(t, outer.InnerError)
}

func TestWrapGenericError(t *testing.T) {

	generic := errors.New("oof. that was bad")
	err := Wrap(generic, "TestEmptyInnerError", "Don't Do This").(SingleError)

	assert.Equal(t, 500, err.Code)
	assert.NotNil(t, err.InnerError)
	assert.Equal(t, "TestEmptyInnerError", err.Location)
	assert.Equal(t, "Don't Do This", err.Message)
	// assert.Equal(t, len(err.Details), 1)

	unwrapped := err.Unwrap()
	assert.Equal(t, "oof. that was bad", unwrapped.Error())
	Report(err)
}

func TestEmptyInnerError(t *testing.T) {

	{
		err := Wrap(nil, "TestEmptyInnerError", "Don't Do This")
		assert.Nil(t, err)
	}

	{
		var innerError error
		outer := Wrap(innerError, "Should Still Be Empty", "Really")
		assert.Nil(t, outer)
	}
}

func TestNotFound(t *testing.T) {

	{
		err := errors.New("regular error")
		require.False(t, NotFound(err))
	}

	{
		err := errors.New("not found")
		require.True(t, NotFound(err))
	}

	{
		err := New(500, "", "")
		require.False(t, NotFound(err))
	}

	{
		err := New(404, "", "")
		require.True(t, NotFound(err))
	}

	{
		e := New(CodeNotFoundError, "Location", "Message")
		assert.Equal(t, CodeNotFoundError, ErrorCode(e))
	}
}

func TestIsNil(t *testing.T) {

	// IsNil has some strange edge cases, so make sure that nobody
	// makes derp panic because they define a strange error type

	{
		var nilPointer *SingleError
		require.True(t, isNil(nilPointer))
	}

	{
		var nilInterface error
		require.True(t, isNil(nilInterface))
	}

	{
		actualError := errors.New("this should not be nil")
		require.False(t, isNil(actualError))
	}

	{
		derpError := New(404, "Code Location", "Error Message")
		require.False(t, isNil(derpError))
	}
}

type weirdErrorType string

func (w weirdErrorType) Error() string {
	return "sure, it's an error"
}

func TestIsNil_WeirdErrorTypes(t *testing.T) {
	{
		require.False(t, isNil(weirdErrorType("")))
	}
}
