package plugins

import (
	"encoding/json"
	"fmt"
)

// JSON prints errors to the system console.
type JSON struct{}

// Report implements the `derp.Plugin` interface, which allows the JSON
// to be called by the derp.Report() method.
func (console JSON) Report(err error) {
	json, _ := json.MarshalIndent(err, "", "\t")
	fmt.Print(string(json))
}
