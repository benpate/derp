package derp

import (
	"errors"
	"testing"
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

func TestReportWrapped(_ *testing.T) {

	inner := newError(codeInternalError, "ouch", "omg")
	middle := Wrap(inner, "whoa", "dude")
	outer := Wrap(middle, "srsly", "bro")

	Report(outer)
}
