package derp

// Reporter wraps the "Report" method, which reports a derp error to an external
// source. Reporters are responsible for handling and swallowing any errors they generate.
type Reporter interface {
	Report(*Error)
}
