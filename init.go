package derp

import "github.com/benpate/derp/plugins"

// Plugins is the list of reporters that are notified whenever Report() is called.
var Plugins ReporterList

// SetPlugins replaces the global reporter list in a single atomic swap.  It is the safe way to
// reconfigure error reporting at RUNTIME: unlike Clear-then-Add, there is never a moment when a
// concurrently reported error reaches an empty list and vanishes.
func SetPlugins(reporters ...Reporter) {
	Plugins.Set(reporters...)
}

// init populates the default plugin list, so that derp reports errors out of the box.
func init() {

	// Start with the JSON reporter as the only item in the list.
	Plugins.Set(plugins.JSON{})
}
