package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiError(t *testing.T) {

	var err error

	err = Append(err, NewBadRequestError("location", "first message"))
	err = Append(err, NewForbiddenError("location", "second message"))

	require.Equal(t, "location: first message", Message(err))
	require.Equal(t, 400, ErrorCode(err))
	require.Equal(t, "location: first message\nlocation: second message\n", err.Error())
}

func TestMultiError_Vanilla(t *testing.T) {
	var err error

	err = Append(err, errors.New(""))
	require.Equal(t, "", Message(err))
}

func TestMultiError_Empty(t *testing.T) {
	err := MultiError{}
	require.Empty(t, Message(err))
	require.Empty(t, ErrorCode(err))
	require.Empty(t, err.Error())
}
