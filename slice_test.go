package derp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlices(t *testing.T) {

	{
		var s []error

		require.Nil(t, s)
	}

	{
		s := []error{}

		require.NotNil(t, s)
		require.Zero(t, len(s))
	}
}
