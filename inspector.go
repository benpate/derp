package derp

// InspectorFunc is a function that can look at any standard error and
// generate a derp.Error out of it.  By default, we include a simple InspectorFunc
// with derp, but using this style allows us to inject a more sophisticated one
// into the library later on.  This more custom InspectorFunc can be used to
// introspect objects that are specific to the application or the specific
// libraries that it uses
type InspectorFunc func(error) *Error

// Inspector is the specific InspectorFunc that is used by this instance of
// the derp library.  By default, this translates generic errors into a derp equivalent,
// but this function be overwritten by the application at runtime (really, init-time).
var Inspector InspectorFunc

func init() {

	Inspector = func(err error) *Error {

		if err == nil {
			return &Error{}
		}

		switch e := err.(type) {

		case *Error:
			return e

		default:
			return &Error{
				Location: "Embedded Error",
				Message:  e.Error(),
				Details:  []interface{}{e},
				Code:     500,
			}
		}
	}
}
