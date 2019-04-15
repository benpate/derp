package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDerp tests basic derp functions (separate from features of a specific reporter)
func TestDerp(t *testing.T) {

	// Create an inner error
	innerError := New("InnerError", "Not Found", CodeNotFoundError, nil, "detail1", "detail2", "detail3")

	assert.Equal(t, innerError.Location, "InnerError")
	assert.Equal(t, innerError.Message, "Not Found")
	assert.Equal(t, innerError.Code, 404)
	assert.Equal(t, innerError.Details[0], "detail1")
	assert.Equal(t, innerError.Details[1], "detail2")
	assert.Equal(t, innerError.Details[2], "detail3")
	assert.Equal(t, innerError.NotFound(), true)

	// Create an outer error
	outerError := New("OuterError", "Inherited", 0, innerError, "other details here")

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
	innerError := New("Location Name", "Error Description", CodeNotFoundError, nil, "details")

	// Verify that the error interface is outputting what we expect.
	assert.Equal(t, innerError.Error(), "Location Name: Error Description")
}

func ExampleNew() error {

	// Mock an error
	thisBreaks := func() error {
		return errors.New("Something failed")
	}

	// Try something that fails
	if err := thisBreaks(); err != nil {

		// Populate a derp.Error with everything you know about the error
		result := New("Example", "Something broke in `thisBreaks`", CodeInternalError, err)

		// Call .Report() to send an error to Ops. This is a system-wide
		// configuration that's set up during initialization.
		result.Report()
	}

	return nil
}
