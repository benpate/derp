package derp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestDerp tests basic derp functions (separate from features of a specific reporter)
func TestDerp(t *testing.T) {

	// Create an inner error
	innerError := New("TestDerp", "Inner Error: Not Found", CodeNotFoundError, nil, "detail1", "detail2", "detail3")

	assert.Equal(t, innerError.Location, "TestDerp")
	assert.Equal(t, innerError.Message, "Inner Error: Not Found")
	assert.Equal(t, innerError.Code, 404)
	assert.Equal(t, innerError.Details[0], "detail1")
	assert.Equal(t, innerError.Details[1], "detail2")
	assert.Equal(t, innerError.Details[2], "detail3")

	// Create an outer error
	outerError := New("TestDerp", "OuterError", 0, innerError, "other details here")

	assert.Equal(t, outerError.Location, "TestDerp")
	assert.Equal(t, outerError.Message, "OuterError")
	assert.Equal(t, outerError.Code, 404) // This is still 404 because we've let the inner error code bubble up
	assert.NotNil(t, outerError.InnerError)
	assert.Equal(t, outerError.Details[0], "other details here")
}
