package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodes(t *testing.T) {

	err := New(123, "whatever", "dude")

	assert.Equal(t, 123, ErrorCode(err))
}

func TestCodeGeneric(t *testing.T) {

	err := errors.New("whatever, dude")

	assert.Equal(t, 500, ErrorCode(err))
}

func TestWithCode(t *testing.T) {
	err := New(123, "whatever", "dude", WithCode(404))
	assert.Equal(t, 404, ErrorCode(err))
}

func TestWithMessage(t *testing.T) {
	err := New(123, "whatever", "dude", WithMessage("message"))
	assert.Equal(t, "message", Message(err))
}
