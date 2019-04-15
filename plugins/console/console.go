package console

import (
	"encoding/json"
	"fmt"

	"github.com/benpate/derp"
)

// Console knows how to report a derp.Error to the system console
type Console struct {
}

// New returns a fully populated Reporter
func New() Console {
	return Console{}
}

// Report sends a derp.Error to the system console. (Great for debugging)
func (console Console) Report(err *derp.Error) {

	json, _ := json.MarshalIndent(err, "", "\t")

	fmt.Print(string(json))

	/*
		sp := spew.NewDefaultConfig()
		sp.ContinueOnMethod = true

		fmt.Print("DERP!! ")
		sp.Dump(err)
	*/
}
