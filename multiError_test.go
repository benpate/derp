package derp

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMultiError_Append(t *testing.T) {

	{
		e := Append(
			Invalid("name", "Name is required"),
			Invalid("email", "Email is required"),
			Invalid("password", "Password is not long enough"),
		)

		require.Equal(t, 3, len(e.Errors))
		require.Equal(t, "name", e.Errors[0].(*ValidationError).Path)
		require.Equal(t, "email", e.Errors[1].(*ValidationError).Path)
		require.Equal(t, "password", e.Errors[2].(*ValidationError).Path)

		Report(e)
	}
}

func TestMultiError_AppendNested(t *testing.T) {

	{
		e := Append(
			Invalid("name", "Name is required"),
			Invalid("email", "Email is required"),
			Append(
				Invalid("password", "Password does not meet complexity requirements"),
				Invalid("confirm_password", "Password entries do not match."),
			),
		)

		require.Equal(t, 4, len(e.Errors))
		require.Equal(t, CodeValidationError, e.ErrorCode())
		require.Equal(t, CodeValidationError, ErrorCode(e))
		require.Equal(t, "name", e.Errors[0].(*ValidationError).Path)
		require.Equal(t, "email", e.Errors[1].(*ValidationError).Path)
		require.Equal(t, "password", e.Errors[2].(*ValidationError).Path)
		require.Equal(t, "confirm_password", e.Errors[3].(*ValidationError).Path)

		Report(e)
	}
}

func TestMultiError_Report(t *testing.T) {

	{
		e := Append(
			Invalid("name", "Name is required"),
			Invalid("email", "Email is required"),
			Invalid("password", "Password is not long enough"),
		)

		Report(e)
	}
}

func TestMultiError_AppendNil(t *testing.T) {

	{
		e := Append(nil, nil, nil)
		require.Nil(t, e)
		require.Zero(t, ErrorCode(e))
	}

	{
		e := &MultiError{}
		require.Zero(t, ErrorCode(e))
	}
}
