package derp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCause_Shallow(t *testing.T) {
	inner := newError(codeForbiddenError, "Inner", "Message")
	outer := Wrap(inner, "Outer", "Message")

	assert.Equal(t, "Inner: Message", RootCause(inner).Error())
	assert.Equal(t, codeForbiddenError, ErrorCode(RootCause(outer)))
	assert.Equal(t, codeForbiddenError, ErrorCode(outer))
}

func TestRootCause_Deep(t *testing.T) {

	e := Wrap(
		Wrap(
			Wrap(
				Wrap(
					newError(123, "Original Location", "Original Message"),
					"Second Location",
					"Second Message",
				),
				"Third Location",
				"Third Message",
			),
			"Fourth Location",
			"Fourth Message",
		),
		"Fifth Location",
		"Fifth Message",
	)

	require.Equal(t, 123, e.(Error).GetErrorCode())
	require.Equal(t, 123, ErrorCode(e))

	rootCause := RootCause(e).(Error)

	require.Equal(t, 123, rootCause.GetErrorCode())
	require.Equal(t, 123, ErrorCode(rootCause))
	require.Equal(t, "Original Location", rootCause.Location)
	require.Equal(t, "Original Message", rootCause.Message)
}
