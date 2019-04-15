package derp

import (
	"fmt"

	"github.com/benpate/derp"
	"github.com/davecgh/go-spew/spew"
)

// Reporter knows how to report a derp.Error to the system console
type Reporter struct {
}

// NewReporter returns a fully populated Reporter
func NewReporter() *Reporter {
	return &Reporter{}
}

// Report sends a derp.Error to the system console. (Great for debugging)
func (reporter *Reporter) Report(err *derp.Error) {

	sp := spew.NewDefaultConfig()
	sp.ContinueOnMethod = true

	fmt.Print("DERP!! ")
	sp.Dump(err)
}
