package derp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestError(t *testing.T) {
	e := New("Location", "Message", 0, nil)
	assert.Equal(t, "Location: Message", e.Error())
}

func TestRootCause(t *testing.T) {
	inner := New("Inner", "Message", CodeForbiddenError, nil)
	outer := New("Outer", "Message", 0, inner)

	assert.Equal(t, "Inner: Message", inner.RootCause().Error())
	assert.Equal(t, CodeForbiddenError, outer.RootCause().Code)
}

func TestReport(t *testing.T) {
	e := New("Location", "Message", 0, nil)
	assert.Equal(t, "Location: Message", e.Error())
}

func TestNotFound(t *testing.T) {
	e := New("Location", "Message", CodeNotFoundError, nil)
	assert.Equal(t, 404, e.Code)
}
