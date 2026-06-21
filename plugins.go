package derp

// Reporter wraps the "Report" method, which reports a derp error to an external
// source. Reporters are responsible for handling and swallowing any errors they generate.
type Reporter interface {
	Report(error)
}

// ReporterList represents an array of reporters, which will be called in succession whenever
// the Report() function is called.
type ReporterList []Reporter

// Add appends a new reporter to this list.  This lets the developer
// configure and append additional reporters during initialization.  It should be called
// during system startup only.
func (list *ReporterList) Add(reporter Reporter) {
	*list = append(*list, reporter)
}

// Clear removes all reporters from this list.  It is useful for
// removing the library default JSON reporter from the list, in the event that
// you don't want to report errors to the console.
func (list *ReporterList) Clear() {
	*list = nil
}
