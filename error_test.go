package derp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	e := New(CodeInternalError, "Location", "Message")
	assert.Equal(t, "Location: Message", e.Error())
}

func TestRootCause(t *testing.T) {
	inner := New(CodeForbiddenError, "Inner", "Message")
	outer := Wrap(inner, "Outer", "Message")

	assert.Equal(t, "Inner: Message", inner.RootCause().Error())
	assert.Equal(t, CodeForbiddenError, outer.RootCause().Code)
}

func TestReport(t *testing.T) {
	e := New(CodeInternalError, "Location", "Message")
	assert.Equal(t, "Location: Message", e.Error())
}

func TestNotFound(t *testing.T) {
	e := New(CodeNotFoundError, "Location", "Message")
	assert.Equal(t, 404, e.Code)
}
