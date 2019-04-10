package derp

import (
	"testing"

	"github.com/benpate/derp"
)

// TestDerp tests basic derp functions (separate from features of a specific reporter)
func TestDerp(t *testing.T) {

	if err := derp.New(); err != nil {

		if err.NotFound() {

		}
	}

	derp.New().Report()
}
