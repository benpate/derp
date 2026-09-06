package derp

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestNewHTTPError_RedactsCredentials confirms a recorded transaction carries no
// credential from either side of the exchange
func TestNewHTTPError_RedactsCredentials(t *testing.T) {

	// A recorded error outlives the transaction: it is wrapped, returned, serialized,
	// and written to a log. A credential that reaches that far has been published.

	request := httptest.NewRequest(http.MethodGet, "https://example.com/api", nil)
	request.Header.Set("Authorization", "Bearer super-secret-token")
	request.Header.Set("Cookie", "session=super-secret-cookie")
	request.Header.Set("Proxy-Authorization", "Basic super-secret-proxy")
	request.Header.Set("X-Api-Key", "super-secret-apikey")
	request.Header.Set("Accept", "application/json")

	response := &http.Response{
		StatusCode: http.StatusUnauthorized,
		Status:     "401 Unauthorized",
		Header:     http.Header{},
	}

	response.Header.Set("Set-Cookie", "session=super-secret-reply")
	response.Header.Set("Content-Type", "application/json")

	err := NewHTTPError(request, response)

	serialized, marshalError := json.Marshal(err)
	require.NoError(t, marshalError)

	require.NotContains(t, string(serialized), "super-secret", "no credential may survive into a recorded error")

	// Headers that carry no credential are still worth having in a report
	require.Equal(t, "application/json", err.Request.Header.Get("Accept"))
	require.Equal(t, "application/json", err.Response.Header.Get("Content-Type"))

	// The redaction is visible rather than silent, so a reader knows a value was sent
	require.Equal(t, redactedValue, err.Request.Header.Get("Authorization"))
	require.Equal(t, redactedValue, err.Response.Header.Get("Set-Cookie"))
}

// TestNewHTTPError_DoesNotMutateTheRequest guards the copy that keeps redaction from
// breaking the transaction it is describing
func TestNewHTTPError_DoesNotMutateTheRequest(t *testing.T) {

	// remote builds an HTTPError from a live *http.Request, and a redirect or a retry
	// reuses it. Redacting in place would strip the credential from the request itself.

	request := httptest.NewRequest(http.MethodGet, "https://example.com/api", nil)
	request.Header.Set("Authorization", "Bearer still-needed")

	response := &http.Response{StatusCode: http.StatusUnauthorized, Header: http.Header{}}
	response.Header.Set("Set-Cookie", "session=still-needed")

	_ = NewHTTPError(request, response)

	require.Equal(t, "Bearer still-needed", request.Header.Get("Authorization"))
	require.Equal(t, "session=still-needed", response.Header.Get("Set-Cookie"))
}

// TestNewHTTPError_RedactionIsCaseInsensitive confirms a header named in any casing
// is still redacted
func TestNewHTTPError_RedactionIsCaseInsensitive(t *testing.T) {

	// http.Header canonicalizes on Set, but a Header built by hand (or decoded) may not,
	// so the lookup canonicalizes rather than trusting the map's keys.
	request := httptest.NewRequest(http.MethodGet, "https://example.com/api", nil)
	request.Header["authorization"] = []string{"Bearer lowercase-secret"}

	err := NewHTTPError(request, nil)

	serialized, marshalError := json.Marshal(err)

	require.NoError(t, marshalError)
	require.NotContains(t, string(serialized), "lowercase-secret")
}

// TestNewHTTPError_AbsentHeadersAreNotInvented confirms redaction does not add a header
// that was never sent
func TestNewHTTPError_AbsentHeadersAreNotInvented(t *testing.T) {

	// Writing "[REDACTED]" for a header nobody sent would report a credential that does
	// not exist, and send the next reader looking for it.
	request := httptest.NewRequest(http.MethodGet, "https://example.com/api", nil)

	err := NewHTTPError(request, nil)

	require.Empty(t, err.Request.Header.Get("Authorization"))
	require.Empty(t, err.Request.Header.Get("Cookie"))
}

// TestNewHTTPError_NilRequestAndResponse confirms the redactor tolerates a transaction
// that never got started
func TestNewHTTPError_NilRequestAndResponse(t *testing.T) {
	err := NewHTTPError(nil, nil)
	require.Equal(t, 0, err.Response.StatusCode)
}
