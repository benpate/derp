package plugins

import (
	"encoding/json"
	"fmt"
)

// JSON prints errors to the system console as indented JSON.
type JSON struct{}

// Report implements the `derp.Reporter` interface, which allows the JSON
// plugin to be called by the derp.Report() method.
func (j JSON) Report(err error) {
	// Per the Reporter contract, reporters swallow their own errors;
	// a marshaling failure here simply prints an empty line.
	bytes, _ := json.MarshalIndent(err, "", "\t")
	fmt.Println(string(bytes))
}
