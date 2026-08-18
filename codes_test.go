package derp

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCodeConstants pins the numeric value of every error code that derp
// defines.  These values are part of the library's public behavior -- callers
// compare them against HTTP status codes -- so changing one is a breaking change.
func TestCodeConstants(t *testing.T) {

	assert.Equal(t, 400, codeBadRequestError)
	assert.Equal(t, 401, codeUnauthorizedError)
	assert.Equal(t, 403, codeForbiddenError)
	assert.Equal(t, 404, codeNotFoundError)
	assert.Equal(t, 409, codeConflictError)
	assert.Equal(t, 410, codeGoneError)
	assert.Equal(t, 418, codeTeapotError)
	assert.Equal(t, 421, codeMisdirectedRequestError)
	assert.Equal(t, 422, codeValidationError)
	assert.Equal(t, 429, codeTooManyRequestsError)
	assert.Equal(t, 500, codeInternalError)
	assert.Equal(t, 501, codeNotImplementedError)
	assert.Equal(t, 502, codeBadGatewayError)
	assert.Equal(t, 524, codeTimeout)
}
