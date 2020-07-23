package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCodes(t *testing.T) {

	err := New(123, "whatever", "dude")

	assert.Equal(t, 123, Code(err))
}

func TestCodeGeneric(t *testing.T) {

	err := errors.New("whatever, dude")

	assert.Equal(t, 500, Code(err))
}
