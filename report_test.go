package derp

import (
	"errors"
	"testing"
)

func TestReport(t *testing.T) {

	err := new(codeNotFoundError, "OMG", "Really Bad")
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

	inner := new(codeInternalError, "ouch", "omg")
	middle := Wrap(inner, "whoa", "dude")
	outer := Wrap(middle, "srsly", "bro")

	Report(outer)
}
