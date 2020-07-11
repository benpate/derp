package derp

import (
	"encoding/json"
	"fmt"
)

// ConsolePlugin prints errors to the system console.
type ConsolePlugin struct{}

// Report implements the `Plugin` interface, which allows the ConsolePlugin
// to be called by the Error.Report() method.
func (consolePlugin ConsolePlugin) Report(err error) {

	if derpError, ok := err.(*Error); ok {
		json, _ := json.MarshalIndent(derpError, "", "\t")
		fmt.Print(string(json))
		return
	}

	// Fall through means this is a regular old error.  Just print what we can.
	fmt.Print(err.Error())
}
