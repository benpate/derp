package derp

// Plugins is the array of objects that are able to report a derp when err.Report() is called.
var Plugins PluginList

func init() {

	// Start with the ConsolePlugin as the only item in the list of plugins.
	Plugins.Add(ConsolePlugin{})
}
