package derp

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

// TestHTTPError_Message_PreferredOverStatus confirms that Error returns the Message, when one
// is set, in place of the status line.
func TestHTTPError_Message_PreferredOverStatus(t *testing.T) {

	err := HTTPError{
		Response: HTTPResponseReport{StatusCode: http.StatusOK, Status: "200 OK"},
		Message:  "Server answered text/html, which cannot be decoded into the requested result",
	}

	require.Equal(t, "Server answered text/html, which cannot be decoded into the requested result", err.Error())
	require.NotContains(t, err.Error(), "200 OK")

	// RULE: The code is untouched.  The message explains the failure; the response still
	// records what the server sent, and GetErrorCode still reports it.
	require.Equal(t, http.StatusOK, err.GetErrorCode())
	require.Equal(t, "200 OK", err.Response.Status)
}

// TestHTTPError_Message_FallsBackToStatus confirms that an HTTPError with no Message still
// explains itself with its status line.
func TestHTTPError_Message_FallsBackToStatus(t *testing.T) {

	err := HTTPError{
		Response: HTTPResponseReport{StatusCode: http.StatusNotFound, Status: "404 Not Found"},
	}

	// RULE: For an ordinary status failure the status line IS the reason, and callers that
	// never set a Message depend on reading it.
	require.Equal(t, "404 Not Found", err.Error())

	// The zero value still works, and still produces no text of its own
	require.Equal(t, "", HTTPError{}.Error())
}

// TestHTTPError_Message_SurvivesToRootMessage confirms that RootMessage returns the Message of
// an HTTPError at the bottom of a wrapped chain.
func TestHTTPError_Message_SurvivesToRootMessage(t *testing.T) {

	httpError := HTTPError{
		Response: HTTPResponseReport{StatusCode: http.StatusOK, Status: "200 OK"},
		Message:  "Server answered text/html, which cannot be decoded into the requested result",
	}

	// RULE: RootMessage prefers the deepest message, and user-facing text and the triage
	// signature are both built on it.
	wrapped := Wrap(Wrap(httpError, "remote.Transaction.decodeResponseBody", "Decoding"), "caller.Load", "Loading")

	require.Equal(t, "Server answered text/html, which cannot be decoded into the requested result", RootMessage(wrapped))

	// PRECONDITION: without a message, the same chain explains itself as the status line.
	// If this ever stops holding, RootMessage has changed and the test above proves nothing.
	bare := Wrap(Wrap(HTTPError{Response: HTTPResponseReport{StatusCode: http.StatusOK, Status: "200 OK"}}, "a", "b"), "c", "d")
	require.Equal(t, "200 OK", RootMessage(bare))
}

// TestHTTPError_Message_DefersToWrappedValue confirms that RootMessage still descends past an
// HTTPError's Message to the inner error it wraps.
func TestHTTPError_Message_DefersToWrappedValue(t *testing.T) {

	httpError := HTTPError{
		Response:     HTTPResponseReport{StatusCode: http.StatusOK, Status: "200 OK"},
		Message:      "Server answered application/json, which could not be decoded",
		WrappedValue: Internal("encoding/json", "unexpected end of JSON input"),
	}

	require.Equal(t, "unexpected end of JSON input", RootMessage(Wrap(httpError, "caller", "Loading")))
}

// TestHTTPError_Message_Serializes confirms that the Message is written to JSON, which is how
// derp-mongo stores errors.
func TestHTTPError_Message_Serializes(t *testing.T) {

	withMessage, err := json.Marshal(HTTPError{Message: "Server answered text/html"})
	require.NoError(t, err)
	require.Contains(t, string(withMessage), `"message":"Server answered text/html"`)

	// RULE: omitempty, so a record without a Message keeps its existing shape
	withoutMessage, err := json.Marshal(HTTPError{})
	require.NoError(t, err)
	require.NotContains(t, string(withoutMessage), `"message"`)
}
