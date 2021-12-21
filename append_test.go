package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAppend(t *testing.T) {

	var result error

	// Starting position is nil
	require.Nil(t, result)

	// Append a nil
	result = Append(result, nil)
	require.Nil(t, result)

	// Append a first error.  Shold just be single error
	result = Append(result, errors.New("omg"))
	require.Equal(t, "omg", result.Error())

	// Append a second error.  Should now be a multi-error
	result = Append(result, errors.New("yet-another-error"))
	require.Equal(t, "omg\nyet-another-error\n", result.Error())

	{
		// Guarantee that we now have a multi-error, with length=2
		multi := result.(MultiError)
		require.Equal(t, 2, len(multi))
	}
}
