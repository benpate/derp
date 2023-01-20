package derp

import (
	"errors"
)

func ExampleNew() {

	// Derp errors work anywhere that you use normal errors.
	// They just contain more information about what actually happened.
	// Here's how to create a new error to report back to a caller
	err := New(404, "Code Location", "Error Message", "additional details here", 12345, map[string]any{})

	// Pluggable error reporting interface can dump errors to the console
	// or anywhere else that you want to send them.
	Report(err)
}

func ExampleWrap() {

	// Derp errors can be nested, containing detailed information
	// about the entire call stack, with specifics about what went
	// wrong at every level

	innerErr := New(404, "Inner Function", "Original Error")

	middleErr := Wrap(innerErr, "Middleware Function", "Error calling 'innerErr'", "parameter", "list", "here")

	outerErr := Wrap(middleErr, "Error in Main Function", "Error calling 'middleErr'", "suspected", "cause", "of", "the", "error")

	Report(outerErr)
}

func ExampleWrap_standardErrors() {

	// Wrap also works with standard library errors
	// so you can add information to errors that are
	// exported by other packages that don't use derp.
	thisBreaks := func() error {
		return errors.New("Something failed")
	}

	// Try something that fails
	if err := thisBreaks(); err != nil {

		// Populate a derp.Error with everything you know about the error
		result := Wrap(err, "Example", "Something broke in `thisBreaks`", "additional details go here")

		// Additional data (such as a custom error code) can be added here.
		SetErrorCode(result, 404)

		// Call .Report() to send an error to Ops. This is a system-wide
		// configuration that's set up during initialization.
		Report(result)
	}
}
