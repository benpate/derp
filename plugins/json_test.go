package plugins

import (
	"errors"
	"testing"
)

// TestJSON_Report verifies that JSON.Report marshals and prints an error
// without panicking. (This package cannot import derp without creating an
// import cycle, so a standard library error is used here.)
func TestJSON_Report(_ *testing.T) {
	JSON{}.Report(errors.New("something went wrong"))
}
