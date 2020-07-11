package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
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
	assert.Equal(t, innerError.NotFound(), true)

	// Create an outer error
	outerError := Wrap(innerError, "OuterError", "Inherited", "other details here")

	assert.Equal(t, outerError.Location, "OuterError")
	assert.Equal(t, outerError.Message, "Inherited")
	assert.Equal(t, outerError.Code, 404) // This is still 404 because we've let the inner error code bubble up
	assert.NotNil(t, outerError.InnerError)
	assert.Equal(t, outerError.Details[0], "other details here")
	assert.Equal(t, outerError.NotFound(), true)

	// Test the RootCause() function
	assert.Equal(t, "InnerError", outerError.RootCause().Location)
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
	err := Wrap(generic, "TestEmptyInnerError", "Don't Do This")

	assert.Equal(t, 500, err.Code)
	assert.Nil(t, err.InnerError)
	assert.Equal(t, "TestEmptyInnerError", err.Location)
	assert.Equal(t, "Don't Do This", err.Message)
	assert.Equal(t, len(err.Details), 1)

	unwrapped := err.Unwrap()
	assert.Equal(t, "oof. that was bad", unwrapped.Error())
}

func TestEmptyInnerError(t *testing.T) {

	err := Wrap(nil, "TestEmptyInnerError", "Don't Do This")

	assert.Equal(t, 500, err.Code)
	assert.Nil(t, err.InnerError)
	assert.Equal(t, "TestEmptyInnerError", err.Location)
	assert.Equal(t, "Don't Do This", err.Message)
	assert.Empty(t, err.Details)

	inner := err.Unwrap()
	assert.Nil(t, inner)
}

func ExampleNew() {

	// Mock an error
	thisBreaks := func() error {
		return errors.New("Something failed")
	}

	// Try something that fails
	if err := thisBreaks(); err != nil {

		// Populate a derp.Error with everything you know about the error
		result := New(CodeInternalError, "Example", "Something broke in `thisBreaks`", err.Error())

		// Call .Report() to send an error to Ops. This is a system-wide
		// configuration that's set up during initialization.
		Report(result)
	}
}
