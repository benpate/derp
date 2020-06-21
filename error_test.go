package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	e := New(CodeInternalError, "Location", "Message")
	assert.Equal(t, "Location: Message", e.Error())
	assert.Equal(t, CodeInternalError, e.ErrorCode())
}

func TestRootCause(t *testing.T) {
	inner := New(CodeForbiddenError, "Inner", "Message")
	outer := Wrap(inner, "Outer", "Message")

	assert.Equal(t, "Inner: Message", inner.RootCause().Error())
	assert.Equal(t, CodeForbiddenError, outer.RootCause().Code)
	assert.Equal(t, CodeForbiddenError, outer.ErrorCode())
}

func TestReport(t *testing.T) {
	e := New(CodeInternalError, "Location", "Message")
	assert.Equal(t, "Location: Message", e.Error())
	assert.Equal(t, CodeInternalError, e.ErrorCode())
}

func TestNotFound(t *testing.T) {
	e := New(CodeNotFoundError, "Location", "Message")
	assert.Equal(t, 404, e.Code)
	assert.Equal(t, CodeNotFoundError, e.ErrorCode())
}

func TestWrapNative(t *testing.T) {

	inner := New(101, "A", "B", "C")
	outer := Wrap(inner, "C", "D")

	assert.Equal(t, outer.Code, 101)
	assert.Equal(t, outer.Location, "C")
	assert.Equal(t, outer.Message, "D")

	innerAgain := outer.Unwrap().(*Error)
	assert.Equal(t, innerAgain.Code, 101)
	assert.Equal(t, innerAgain.Location, "A")
	assert.Equal(t, innerAgain.Message, "B")
}

func TestWrapGeneric(t *testing.T) {

	inner := errors.New("omg it works")
	outer := Wrap(inner, "C", "D")

	assert.Equal(t, outer.Code, 500)
	assert.Equal(t, outer.Location, "C")
	assert.Equal(t, outer.Message, "D")

	innerAgain := outer.Unwrap()
	assert.Equal(t, innerAgain.Error(), "omg it works")
}
