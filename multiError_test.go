package derp

import (
	"errors"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type testCustomCodeError int

func (t testCustomCodeError) Error() string {
	return "Custom Code: " + strconv.Itoa(int(t))
}

func (t testCustomCodeError) ErrorCode() int {
	return int(t)
}

func TestMultiError_Append(t *testing.T) {

	{
		e := Append(
			errors.New("first error here"),
			errors.New("second error here"),
			errors.New("third error here"),
		)

		require.Equal(t, 3, len(*e))
		require.Equal(t, "first error here", (*e)[0].Error())
		require.Equal(t, "second error here", (*e)[1].Error())
		require.Equal(t, "third error here", (*e)[2].Error())

		Report(e)
	}
}

func TestMultiError_AppendNested(t *testing.T) {

	{
		e := Append(
			errors.New("first error here"),
			errors.New("second error here"),
			Append(
				errors.New("first nested error here"),
				errors.New("second nested error here"),
			),
		)

		require.Equal(t, 4, len(*e))
		require.Equal(t, 500, e.ErrorCode())
		require.Equal(t, 500, ErrorCode(e))
		require.Equal(t, "first error here", (*e)[0].Error())
		require.Equal(t, "second error here", (*e)[1].Error())
		require.Equal(t, "first nested error here", (*e)[2].Error())
		require.Equal(t, "second nested error here", (*e)[3].Error())

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

func TestMultiError_Code(t *testing.T) {

	{
		e := Append(
			testCustomCodeError(101),
			testCustomCodeError(202),
			testCustomCodeError(303),
		)
		assert.Equal(t, 101, ErrorCode(e))
	}

	{
		e := Append(
			errors.New("this has no error code"),
			testCustomCodeError(202),
			testCustomCodeError(303),
		)
		assert.Equal(t, 202, ErrorCode(e))
	}

	{
		e := Append(
			errors.New("this has no error code"),
			errors.New("this has no error code"),
			errors.New("this has no error code"),
		)

		assert.Equal(t, 500, ErrorCode(e))
	}
}
