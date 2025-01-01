package derp

import "reflect"

/******************************************
 * Other Helpers
 ******************************************/

// isNil performs a robust nil check on an error interface
// Shout out to: https://medium.com/@mangatmodi/go-check-nil-interface-the-right-way-d142776edef1
func isNil(i error) bool {

	if i == nil {
		return true
	}

	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Array, reflect.Slice, reflect.Chan, reflect.Map:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}
