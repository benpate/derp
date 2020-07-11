package derp

import (
	"encoding/json"
	"fmt"
)

// ConsolePlugin prints errors to the system console.
type ConsolePlugin struct{}

// Report implements the `Plugin` interface, which allows the ConsolePlugin
// to be called by the Error.Report() method.
func (consolePlugin ConsolePlugin) Report(err *Error) {

	json, _ := json.MarshalIndent(err, "", "\t")
	fmt.Print(string(json))
	return
}
