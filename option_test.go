package derp

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestOption(t *testing.T) {

	var f any = WithNotFound()

	if _, ok := f.(Option); ok {
		t.Log("f is an Option")
	} else {
		t.Error("f is not an Option")
	}
}

func TestOption_New(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithCode(codeInternalError))
	assert.Equal(t, codeInternalError, ErrorCode(e))
}

func TestOption_Wrap(t *testing.T) {

	e := errors.New("wrapped error")
	wrapped := Wrap(e, "Location", "Message", WithCode(codeInternalError))
	assert.Equal(t, codeInternalError, ErrorCode(wrapped))
}

func TestOption_WithBadRequest(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithBadRequest())
	assert.Equal(t, codeBadRequestError, e.Code)
}

func TestOption_WithForbidden(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithForbidden())
	assert.Equal(t, codeForbiddenError, e.Code)
}

func TestOption_WithInternalError(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithInternalError())
	assert.Equal(t, codeInternalError, e.Code)
}

func TestOption_WithNotFound(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithNotFound())
	assert.Equal(t, codeNotFoundError, e.Code)
}

func TestOption_WithLocation(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithLocation("New Location"))
	assert.Equal(t, "New Location", e.Location)
}

func TestOption_WithMessage(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithMessage("New Message"))
	assert.Equal(t, "New Message", e.Message)
}

func TestOption_WithWrappedValue(t *testing.T) {
	e := newError(codeNotFoundError, "Location", "Message", WithWrappedValue(errors.New("wrapped error")))
	assert.Equal(t, "wrapped error", e.WrappedValue.Error())
}
