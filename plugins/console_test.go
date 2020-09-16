package plugins

import (
	"strconv"
	"testing"
)

type testError struct {
	Code     int
	Location string
	Message  string
}

func (t testError) Error() string {
	return strconv.Itoa(t.Code) + ": " + t.Location + ": " + t.Message
}

// TestConsole uses the default console plugin to report errors directly
// to the console.
func TestConsole(t *testing.T) {

	// Create a new error
	err1 := testError{
		Code:     500,
		Location: "TestConsole",
		Message:  "Error 1",
	}

	console := Console{}

	console.Report(err1)
}
