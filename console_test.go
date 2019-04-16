package derp

import (
	"testing"
)

// TestConsole uses the default console plugin to report errors directly
// to the console.
func TestConsole(t *testing.T) {

	// Create a new error
	err1 := New(CodeInternalError, "TestConsole", "Error 1")

	err2 := Wrap(err1, "TestConsole", "Error 2")

	err3 := Wrap(err2, "TestConsole", "Error 3")

	err3.Report()

	// t.Fail() // Uncomment this to show output without using -v
}
