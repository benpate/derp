package derp

// Report takes ANY error (hopefully a derp error) and attempts to report it
// via all configured error reporting mechanisms.
func Report(err error) {

	// If this is a natural derp error, then report it through all reporting mechanisms
	if derpError, ok := err.(*Error); ok {
		for _, plugin := range Plugins {
			plugin.Report(derpError)
		}

		return
	}

	// Fall through to here means it is NOT a derp error.  Wrap the original in a
	// new derp, and then report.
	Report(Wrap(err, "", ""))
}
