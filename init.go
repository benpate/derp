package derp

import "github.com/benpate/derp/plugins"

// Plugins is the list of reporters that are notified whenever Report() is called.
var Plugins ReporterList = make([]Reporter, 0)

func init() {

	// Start with the JSON reporter as the only item in the list.
	Plugins.Add(plugins.JSON{})
}
