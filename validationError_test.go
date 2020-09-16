package derp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestValidationError(t *testing.T) {

	e := Invalid("name", "name is required").(*ValidationError)

	require.Equal(t, "name", e.Path)
	require.Equal(t, "name is required", e.Message)
	require.Equal(t, "name is required", e.Error())
	require.Equal(t, CodeValidationError, e.ErrorCode())
}
