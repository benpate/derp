package derp

// Report takes ANY error (hopefully a derp error) and attempts to report it
// via all configured error reporting mechanisms.
func Report(err error) error {

	// If the error is nil, then there's nothing to do.
	if err == nil {
		return nil
	}

	// If this is a natural derp error, then report it through all reporting mechanisms
	if derpError, ok := err.(*Error); ok {
		for _, plugin := range Plugins {
			plugin.Report(derpError)
		}

		return derpError
	}

	// Fall through to here means it is NOT a derp error.  Wrap the original in a
	// new derp, and then report.
	return Report(Wrap(err, "derp.Report", "Reporting Generic Error"))
}
