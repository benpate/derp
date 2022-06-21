package derp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNil(t *testing.T) {

	var err error

	require.Zero(t, ErrorCode(err))
	require.Empty(t, Message(err))
	require.False(t, NotFound(err))
	SetErrorCode(err, 500) // this should not break
}
