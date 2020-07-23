package derp

import (
	"errors"
	"testing"
)

func TestReport(t *testing.T) {

	err := New(42, "OMG", "Really Bad")
	Report(err)
}

func TestReportGeneric(t *testing.T) {

	err := errors.New("OMG Really Bad")
	Report(err)
}

func TestReportNil(t *testing.T) {
	Report(nil)
}

func TestReportWrapped(t *testing.T) {

	inner := New(500, "ouch", "omg")
	middle := Wrap(inner, "whoa", "dude")
	outer := Wrap(middle, "srsly", "bro")

	Report(outer)
}
