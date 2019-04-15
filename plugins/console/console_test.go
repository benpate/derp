package console

import (
	"testing"

	"github.com/benpate/derp"
)

func TestConsole(t *testing.T) {

	// Add a reporter plugin to the derp package.  In regular usage,
	// this would only happen once upon application initialization.
	derp.Connect(New())

	// Create a new error
	err1 := derp.New("TestConsole", "Error 1", derp.CodeInternalError, nil)

	err2 := derp.New("TestConsole", "Error 2", derp.CodeNotFoundError, err1)

	err3 := derp.New("TestConsole", "Error 3", derp.CodeInternalError, err2)

	err3.Report()

	// t.Fail() // Uncomment this to show output without using -v
}
