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
