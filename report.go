package derp

// Report takes ANY error (hopefully a derp error) and attempts to report it
// via all configured error reporting mechanisms.
func Report(err error) {

	// If the error is NOT nil, then send "Report" to each installed plugin.
	if NotNil(err) {
		for _, plugin := range Plugins {
			plugin.Report(err)
		}
	}
}

// ReportFunc executes the provided function and reports any error that occurs.
// This is useful with `defer` statements, preventing the underlying function
// from being executed before the deferred call.
func ReportFunc(fn func() error) {
	Report(fn())
}

// ReportAndReturn reports an error to the logger
// and also returns it to the caller.
func ReportAndReturn(err error) error {
	Report(err)
	return err
}
