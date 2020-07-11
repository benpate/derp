package derp

// Report takes ANY error (hopefully a derp error) and attempts to report it
// via all configured error reporting mechanisms.
func Report(err error) {

	for _, plugin := range Plugins {
		plugin.Report(err)
	}
}
