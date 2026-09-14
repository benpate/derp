package derp

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// bodyOf builds a response carrying the supplied body, as a live ReadCloser
func bodyOf(body string) *http.Response {
	return &http.Response{
		StatusCode: http.StatusBadRequest,
		Status:     "400 Bad Request",
		Header:     http.Header{},
		Body:       io.NopCloser(strings.NewReader(body)),
	}
}

// TestNewHTTPError_CapturesResponseBody confirms the server's own explanation reaches
// the report
func TestNewHTTPError_CapturesResponseBody(t *testing.T) {

	// A status line says a request failed; the body says why. Recording the first
	// without the second is what makes a 401 in a log unactionable.
	response := bodyOf(`{"error":"invalid_grant","description":"code expired"}`)

	err := NewHTTPError(nil, response)

	require.Equal(t, `{"error":"invalid_grant","description":"code expired"}`, err.Response.Body)
	require.False(t, err.Response.BodyTruncated)
}

// TestNewHTTPError_LeavesBodyReadable is the guard that keeps this feature from
// breaking the transaction it describes
func TestNewHTTPError_LeavesBodyReadable(t *testing.T) {

	// remote builds the report and THEN decodes the body into its failure object.
	// Consuming the stream here would hand that decoder an empty document.
	response := bodyOf("the whole body, every byte of it")

	_ = NewHTTPError(nil, response)

	remaining, readError := io.ReadAll(response.Body)

	require.NoError(t, readError)
	require.Equal(t, "the whole body, every byte of it", string(remaining))
}

// TestNewHTTPError_LeavesLongBodyReadable confirms the full body survives even when
// only part of it was captured
func TestNewHTTPError_LeavesLongBodyReadable(t *testing.T) {

	// The cap bounds what the REPORT carries, not what the caller receives. Truncating
	// the caller's copy would corrupt a response that is merely large.
	full := strings.Repeat("A", maxResponseBodyBytes*2)
	response := bodyOf(full)

	err := NewHTTPError(nil, response)

	require.True(t, err.Response.BodyTruncated)
	require.Len(t, err.Response.Body, maxResponseBodyBytes)

	remaining, readError := io.ReadAll(response.Body)

	require.NoError(t, readError)
	require.Equal(t, full, string(remaining), "the caller still gets the entire body")
}

// TestNewHTTPError_BodyCapIsFourKilobytes pins the cap itself
func TestNewHTTPError_BodyCapIsFourKilobytes(t *testing.T) {
	require.Equal(t, 4096, maxResponseBodyBytes)
}

// TestNewHTTPError_BodyExactlyAtCapIsNotTruncated confirms the off-by-one at the
// boundary reports honestly
func TestNewHTTPError_BodyExactlyAtCapIsNotTruncated(t *testing.T) {

	// A body that ends exactly at the cap is complete. Reporting it as truncated would
	// send a reader looking for bytes that do not exist.
	full := strings.Repeat("B", maxResponseBodyBytes)
	response := bodyOf(full)

	err := NewHTTPError(nil, response)

	require.Equal(t, full, err.Response.Body)
	require.False(t, err.Response.BodyTruncated)
}

// TestNewHTTPError_ClosesThroughToTheRealBody confirms the connection is still released
func TestNewHTTPError_ClosesThroughToTheRealBody(t *testing.T) {

	// The replacement body wraps the original. If Close stopped at the wrapper, the
	// underlying connection would never be returned to the pool.
	tracker := &closeTracker{Reader: strings.NewReader("body")}
	response := &http.Response{StatusCode: http.StatusBadRequest, Header: http.Header{}, Body: tracker}

	_ = NewHTTPError(nil, response)

	require.NoError(t, response.Body.Close())
	require.True(t, tracker.closed, "Close must reach the real body")
}

// TestNewHTTPError_EmptyBody confirms an absent body adds nothing to the report
func TestNewHTTPError_EmptyBody(t *testing.T) {

	// `omitempty` keeps the field out of the JSON entirely, so a reader is never shown
	// an empty string that looks like a body the server actually sent.
	response := bodyOf("")

	err := NewHTTPError(nil, response)

	require.Empty(t, err.Response.Body)

	serialized, marshalError := json.Marshal(err)

	require.NoError(t, marshalError)
	require.NotContains(t, string(serialized), "body")
}

// TestNewHTTPError_NilBody confirms a response that never had a body is tolerated
func TestNewHTTPError_NilBody(t *testing.T) {

	// A response built by hand (and every test in this package before now) has a nil
	// Body, which must not panic.
	response := &http.Response{StatusCode: http.StatusBadRequest, Header: http.Header{}}

	require.NotPanics(t, func() {
		err := NewHTTPError(nil, response)
		require.Empty(t, err.Response.Body)
	})
}

// TestNewHTTPError_FailingBodyIsLeftAlone confirms a body that cannot be read is not
// replaced by an empty one
func TestNewHTTPError_FailingBodyIsLeftAlone(t *testing.T) {

	// remote closes the body on one of its error paths, and reading a closed body
	// yields nothing. Swapping in a replacement there would hide the original's error
	// from the code that still owns it.
	failing := &failingBody{}
	response := &http.Response{StatusCode: http.StatusBadGateway, Header: http.Header{}, Body: failing}

	err := NewHTTPError(nil, response)

	require.Empty(t, err.Response.Body)
	require.Equal(t, io.ReadCloser(failing), response.Body, "the original body is untouched")
}

// closeTracker records whether Close reached the underlying body
type closeTracker struct {
	io.Reader
	closed bool
}

func (tracker *closeTracker) Close() error {
	tracker.closed = true
	return nil
}

// failingBody models a body that errors on every read, such as one already closed
type failingBody struct{}

func (body *failingBody) Read([]byte) (int, error) {
	return 0, io.ErrUnexpectedEOF
}

func (body *failingBody) Close() error {
	return nil
}

// TestNewHTTPError_MatchesRemoteSequence reproduces the order remote actually uses
func TestNewHTTPError_MatchesRemoteSequence(t *testing.T) {

	// remote.Transaction reads the body into a buffer and swaps in a re-readable
	// NopCloser, THEN builds the report on a non-2xx, THEN lets its caller read the
	// body again. Capturing has to survive being in the middle of that.
	payload := `{"error":"mismatched signature"}`
	response := bodyOf(payload)

	// 1. remote reads and replaces the body (transaction-response.go)
	buffered, readError := io.ReadAll(response.Body)
	require.NoError(t, readError)
	response.Body = io.NopCloser(strings.NewReader(string(buffered)))

	// 2. remote builds the report on a failing status (processResponse)
	err := NewHTTPError(nil, response)
	require.Equal(t, payload, err.Response.Body, "the reason for the failure is recorded")

	// 3. the caller reads the body afterward, and must still see all of it
	afterwards, secondError := io.ReadAll(response.Body)
	require.NoError(t, secondError)
	require.Equal(t, payload, string(afterwards))
}
