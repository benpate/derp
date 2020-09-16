package derp

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestRootCause_Shallow(t *testing.T) {
	inner := New(CodeForbiddenError, "Inner", "Message")
	outer := Wrap(inner, "Outer", "Message")

	assert.Equal(t, "Inner: Message", RootCause(inner).Error())
	assert.Equal(t, CodeForbiddenError, ErrorCode(RootCause(outer)))
	assert.Equal(t, CodeForbiddenError, ErrorCode(outer))
}

func TestRootCause_Deep(t *testing.T) {

	e := Wrap(
		Wrap(
			Wrap(
				Wrap(
					New(123, "Original Location", "Original Message"),
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

	require.Equal(t, 123, e.ErrorCode())
	require.Equal(t, 123, ErrorCode(e))

	rootCause := RootCause(e).(*SingleError)

	require.Equal(t, 123, rootCause.ErrorCode())
	require.Equal(t, 123, ErrorCode(rootCause))
	require.Equal(t, "Original Location", rootCause.Location)
	require.Equal(t, "Original Message", rootCause.Message)
}
