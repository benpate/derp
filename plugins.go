package derp

// Plugin wraps the "Report" method, which reports a derp error to an external
// source. Reporters are responsible for handling and swallowing any errors they generate.
type Plugin interface {
	Report(error)
}

// PluginList represents an array of plugins, which will be called in succession whenever
// the Error.Report() function is called.
type PluginList []Plugin

// Add appends a new plugin to this list.  This lets the developer
// configure and append additional plugins during initialization.  It should be called
// during system startup only.
func (list *PluginList) Add(plugin Plugin) {
	*list = append(*list, plugin)
}

// Clear removes all plugins from this list.  It is useful for
// removing the library default Console() from the list of plugins, in the event that
// you don't want to report errors to the console.
func (list *PluginList) Clear() {
	*list = nil
}
