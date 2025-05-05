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
	innerError := NewNotFoundError("WrappedValue", "Not Found", "detail1", "detail2", "detail3")

	assert.Equal(t, innerError.Location, "WrappedValue")
	assert.Equal(t, innerError.Message, "Not Found")
	assert.Equal(t, innerError.Code, 404)
	assert.Equal(t, innerError.Details[0], "detail1")
	assert.Equal(t, innerError.Details[1], "detail2")
	assert.Equal(t, innerError.Details[2], "detail3")
	assert.Equal(t, NotFound(innerError), true)

	// Create an outer error
	outerError := Wrap(innerError, "OuterError", "Inherited", "other details here").(Error)

	assert.Equal(t, outerError.Location, "OuterError")
	assert.Equal(t, outerError.Message, "Inherited")
	assert.Equal(t, outerError.Code, 404) // This is still 404 because we've let the inner error code bubble up
	assert.NotNil(t, outerError.WrappedValue)
	assert.Equal(t, outerError.Details[0], "other details here")
	assert.Equal(t, NotFound(outerError), true)

	// Test the RootCause() function
	assert.Equal(t, "WrappedValue", RootCause(outerError).(Error).Location)
}

func TestConvenienceFns(t *testing.T) {

	badRequest := NewBadRequestError("location", "description")
	require.Equal(t, codeBadRequestError, ErrorCode(badRequest))

	forbidden := NewForbiddenError("location", "description")
	require.Equal(t, codeForbiddenError, ErrorCode(forbidden))

	internal := NewInternalError("location", "description")
	require.Equal(t, codeInternalError, ErrorCode(internal))

	notFound := NewNotFoundError("location", "description")
	require.Equal(t, codeNotFoundError, ErrorCode(notFound))

	unauthorized := NewUnauthorizedError("location", "description")
	require.Equal(t, codeUnauthorizedError, ErrorCode(unauthorized))

	invalid := NewValidationError("location", "description")
	require.Equal(t, codeValidationError, ErrorCode(invalid))

	teapot := NewTeapotError("location", "description")
	require.Equal(t, codeTeapotError, ErrorCode(teapot))

	misdirected := NewMisdirectedRequestError("location", "description")
	require.Equal(t, codeMisdirectedRequestError, ErrorCode(misdirected))

	notImplemented := NewNotImplementedError("location", "description")
	require.Equal(t, codeNotImplementedError, ErrorCode(notImplemented))
}

func TestMessage(t *testing.T) {

	require.Equal(t, "", Message(nil))

	derp := NewNotFoundError("location", "description")
	require.Equal(t, "description", Message(derp))

	standard := errors.New("this is a standard error")
	require.Equal(t, "this is a standard error", Message(standard))
}

func TestErrorInterface(t *testing.T) {

	// Create an error
	innerError := NewNotFoundError("Location Name", "Error Description", "details")

	// Verify that the error interface is outputting what we expect.
	assert.Equal(t, innerError.Error(), "Location Name: Error Description")
}

func TestStandardError(t *testing.T) {

	// Testing how derp handles an error from the standard library
	err := errors.New("This is a standard error")

	// Wrap it the stdlib error in a derp.  This means: 1) assigning an error code, and 2) making the original error message a property of the derp.Error.
	outer := NewInternalError("TestStandardError", "Encapsulating Error", err.Error())

	assert.Equal(t, "TestStandardError", outer.Location)
	assert.Equal(t, "Encapsulating Error", outer.Message)
	assert.Equal(t, 500, outer.Code)
	assert.Equal(t, 1, len(outer.Details))
	assert.Nil(t, outer.WrappedValue)
}

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

func TestEmptyWrappedValue(t *testing.T) {

	{
		err := Wrap(nil, "TestEmptyWrappedValue", "Don't Do This")
		assert.Nil(t, err)
	}

	{
		var innerError error
		outer := Wrap(innerError, "Should Still Be Empty", "Really")
		assert.Nil(t, outer)
	}
}

func TestNotFound(t *testing.T) {

	require.False(t, NotFound(nil))

	{
		err := errors.New("regular error")
		require.False(t, NotFound(err))
	}

	{
		err := errors.New("not found")
		require.True(t, NotFound(err))
	}

	{
		err := new(500, "", "")
		require.False(t, NotFound(err))
	}

	{
		err := new(404, "", "")
		require.True(t, NotFound(err))
	}

	{
		e := NewNotFoundError("Location", "Message")
		assert.Equal(t, 404, ErrorCode(e))
	}
}

func TestIsNil(t *testing.T) {

	// IsNil has some strange edge cases, so make sure that nobody
	// makes derp panic because they define a strange error type

	{
		var nilPointer *Error
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
		derpError := new(404, "Code Location", "Error Message")
		require.False(t, isNil(derpError))
	}
}

func TestNilOrNotFound(t *testing.T) {

	require.True(t, NilOrNotFound(nil))

	{
		err := errors.New("not found")
		require.True(t, NilOrNotFound(err))
	}

	{
		err := new(500, "", "")
		require.False(t, NilOrNotFound(err))
	}

	{
		err := new(404, "", "")
		require.True(t, NilOrNotFound(err))
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

func TestNilErrorCode(t *testing.T) {
	require.Equal(t, 0, ErrorCode(nil))
}

func TestReportAndReturn(t *testing.T) {

	{
		err := errors.New("regular error")
		require.Equal(t, err, ReportAndReturn(err))
	}

	{
		err := new(404, "Location", "Message")
		require.Equal(t, err, ReportAndReturn(err))
	}
}

func TestIsInformational(t *testing.T) {
	{
		e := new(99, "location", "message")
		require.False(t, IsInformational(e))
	}
	{
		e := new(100, "Location", "Message")
		require.True(t, IsInformational(e))
	}
	{
		e := new(199, "Location", "Message")
		require.True(t, IsInformational(e))
	}
	{
		e := new(200, "Location", "Message")
		require.False(t, IsInformational(e))
	}
}

func TestIsSuccess(t *testing.T) {
	{
		e := new(199, "location", "message")
		require.False(t, IsSuccess(e))
	}
	{
		e := new(200, "Location", "Message")
		require.True(t, IsSuccess(e))
	}
	{
		e := new(299, "Location", "Message")
		require.True(t, IsSuccess(e))
	}
	{
		e := new(300, "Location", "Message")
		require.False(t, IsSuccess(e))
	}
}

func TestIsRedirection(t *testing.T) {
	{
		e := new(299, "location", "message")
		require.False(t, IsRedirection(e))
	}
	{
		e := new(300, "Location", "Message")
		require.True(t, IsRedirection(e))
	}
	{
		e := new(399, "Location", "Message")
		require.True(t, IsRedirection(e))
	}
	{
		e := new(400, "Location", "Message")
		require.False(t, IsRedirection(e))
	}
}

func TestIsClientError(t *testing.T) {
	{
		e := new(399, "location", "message")
		require.False(t, IsClientError(e))
	}
	{
		e := new(400, "Location", "Message")
		require.True(t, IsClientError(e))
	}
	{
		e := new(499, "Location", "Message")
		require.True(t, IsClientError(e))
	}
	{
		e := new(500, "Location", "Message")
		require.False(t, IsClientError(e))
	}
}

func TestIsServerError(t *testing.T) {
	{
		e := new(499, "location", "message")
		require.False(t, IsServerError(e))
	}
	{
		e := new(500, "Location", "Message")
		require.True(t, IsServerError(e))
	}
	{
		e := new(599, "Location", "Message")
		require.True(t, IsServerError(e))
	}
	{
		e := new(600, "Location", "Message")
		require.False(t, IsServerError(e))
	}
}
