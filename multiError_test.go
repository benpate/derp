package derp

import (
	"errors"
	"testing"

	"github.com/davecgh/go-spew/spew"
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

func TestMultiError_Append(t *testing.T) {

	var err MultiError
	spew.Config.DisableMethods = true

	err.Append(NewBadRequestError("location", "first"))
	err.Append(NewForbiddenError("location", "second"))
	err.Append(NewInternalError("location", "third"))

	// require.Equal(t, 3, err.Length())

	err.AddPrefixes("prefix.")
	require.Equal(t, "location: prefix.first", Message(err[0]))
	require.Equal(t, "location: prefix.second", Message(err[1]))
	require.Equal(t, "location: prefix.third", Message(err[2]))
}
