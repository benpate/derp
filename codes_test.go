package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCodes(t *testing.T) {

	err := New(123, "whatever", "dude")

	assert.Equal(t, 123, ErrorCode(err))
}

func TestCodeGeneric(t *testing.T) {

	err := errors.New("whatever, dude")

	assert.Equal(t, 500, ErrorCode(err))
}

func TestSetCode(t *testing.T) {

	{
		var err *SingleError
		SetErrorCode(err, 404)
		require.Equal(t, 0, ErrorCode(err))
	}

	{
		var err error
		SetErrorCode(err, 404)
		require.Equal(t, 0, ErrorCode(err))
	}

	{
		err := &SingleError{}
		SetErrorCode(err, 404)
		require.Equal(t, 404, ErrorCode(err))
	}
}
