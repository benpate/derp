package derp

import (
	"errors"
	"fmt"
)

func ExampleNotFound() {

	// Derp errors work anywhere that you use normal errors.
	// They just contain more information about what actually happened.
	// Here's how to create a new error to report back to a caller.
	err := NotFound("Code Location", "Error Message", "additional details here", 12345, map[string]any{})

	// The error carries a code, location, and message that callers can inspect.
	fmt.Println(ErrorCode(err))
	fmt.Println(Location(err))
	fmt.Println(Message(err))
	// Output:
	// 404
	// Code Location
	// Error Message
}

func ExampleWrap() {

	// Derp errors can be nested, containing detailed information
	// about the entire call stack, with specifics about what went
	// wrong at every level.
	innerErr := NotFound("Inner Function", "Original Error")

	middleErr := Wrap(innerErr, "Middleware Function", "Error calling 'innerErr'", "parameter", "list", "here")

	outerErr := Wrap(middleErr, "Outer Function", "Error calling 'middleErr'", "suspected", "cause", "of", "the", "error")

	// The wrapped code propagates outward, while the root location and message
	// remain reachable through the chain.
	fmt.Println(ErrorCode(outerErr))
	fmt.Println(RootLocation(outerErr))
	fmt.Println(RootMessage(outerErr))
	// Output:
	// 404
	// Inner Function
	// Original Error
}

func ExampleWrap_standardErrors() {

	// Wrap also works with standard library errors
	// so you can add information to errors that are
	// exported by other packages that don't use derp.
	thisBreaks := func() error {
		return errors.New("something failed")
	}

	// Try something that fails.
	err := thisBreaks()

	// Populate a derp.Error with everything you know about the error.
	result := Wrap(err, "Example", "Something broke in `thisBreaks`", WithCode(404), "additional details go here")

	// The standard error is preserved in the chain and matchable via errors.Is.
	fmt.Println(ErrorCode(result))
	fmt.Println(errors.Is(result, err))
	// Output:
	// 404
	// true
}
