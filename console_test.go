package derp

import (
	"testing"
)

// TestConsole uses the default console plugin to report errors directly
// to the console.
func TestConsole(t *testing.T) {

	// Create a new error
	err1 := New("TestConsole", "Error 1", CodeInternalError, nil)

	err2 := New("TestConsole", "Error 2", CodeNotFoundError, err1)

	err3 := New("TestConsole", "Error 3", CodeInternalError, err2)

	err3.Report()

	// t.Fail() // Uncomment this to show output without using -v
}
