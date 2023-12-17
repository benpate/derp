package derp

import "testing"

func TestOption(t *testing.T) {

	var f any = WithNotFound()

	if _, ok := f.(Option); ok {
		t.Log("f is an Option")
	} else {
		t.Error("f is not an Option")
	}
}
