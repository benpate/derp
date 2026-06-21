package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReport(_ *testing.T) {

	err := newError(codeNotFoundError, "OMG", "Really Bad")
	Report(err)
}

func TestReportGeneric(_ *testing.T) {

	err := errors.New("OMG Really Bad")
	Report(err)
}

func TestReportNil(_ *testing.T) {
	Report(nil)
}

func TestReportFunc(t *testing.T) {

	// Swap in a counting plugin, restoring the global list afterwards so
	// other tests are not affected.
	original := Plugins
	t.Cleanup(func() { Plugins = original })

	counter := &countingPlugin{}
	Plugins = ReporterList{counter}

	// ReportFunc should report the error returned by the function
	ReportFunc(func() error {
		return NotFound("location", "message")
	})
	require.Equal(t, 1, counter.count)

	// ReportFunc should NOT report when the function returns nil
	ReportFunc(func() error {
		return nil
	})
	require.Equal(t, 1, counter.count)
}

func TestReportWrapped(_ *testing.T) {

	inner := newError(codeInternalError, "ouch", "omg")
	middle := Wrap(inner, "whoa", "dude")
	outer := Wrap(middle, "srsly", "bro")

	Report(outer)
}
