package derp

import (
	"bytes"
	"io"
	"net/http"
	"net/textproto"
	"strconv"
	"time"
)

// HTTPError wraps a standard derp.Error, including additional data about a failed HTTP transaction
type HTTPError struct {
	Request  HTTPRequestReport  `json:"request"`  // Details of the HTTP request that failed
	Response HTTPResponseReport `json:"response"` // Details of the HTTP response that was returned

	WrappedValue error `json:"innerError,omitempty"` // An underlying error object used to identify the root cause of this error.
}

// NewHTTPError creates a new HTTPError object from the given request and response
func NewHTTPError(request *http.Request, response *http.Response) HTTPError {

	result := HTTPError{}

	if request != nil {

		result.Request = HTTPRequestReport{
			Method: request.Method,
			Header: RedactHeader(request.Header),
		}

		// RULE: A Request is not guaranteed to have a URL, which cannot be stringified when nil.
		if request.URL != nil {
			result.Request.URL = request.URL.String()
		}
	}

	if response != nil {

		body, truncated := captureResponseBody(response)

		result.Response = HTTPResponseReport{
			StatusCode:    response.StatusCode,
			Status:        response.Status,
			Header:        RedactHeader(response.Header),
			Body:          body,
			BodyTruncated: truncated,
		}
	}

	return result
}

// maxResponseBodyBytes caps how much of a response body is copied into an error report
const maxResponseBodyBytes = 4096

// captureResponseBody copies the first maxResponseBodyBytes of the response body for
// the report, reporting whether more was left behind.
func captureResponseBody(response *http.Response) (string, bool) {

	if response == nil || response.Body == nil {
		return "", false
	}

	original := response.Body

	// RULE: read one byte past the cap. That is what distinguishes a body that ends at
	// exactly the cap from one that was cut off, so the report can say which it was
	// instead of ending mid-sentence and leaving the reader to guess.
	//
	// RULE: a read error is not reported. This is already the error path, the body is
	// best-effort evidence, and the bytes that did arrive are still the most useful
	// thing in the report.
	captured, _ := io.ReadAll(io.LimitReader(original, maxResponseBodyBytes+1))

	// Nothing was consumed, so there is nothing to put back. Leaving the body untouched
	// also preserves whatever error it is holding for its owner to find.
	if len(captured) == 0 {
		return "", false
	}

	// RULE: restore before returning. The body belongs to a live response whose owner
	// may still decode it, and Close must still reach the real body or the underlying
	// connection is never released.
	response.Body = replayBody{
		Reader: io.MultiReader(bytes.NewReader(captured), original),
		Closer: original,
	}

	if len(captured) > maxResponseBodyBytes {
		return string(captured[:maxResponseBodyBytes]), true
	}

	return string(captured), false
}

// replayBody re-serves bytes captured for an error report, keeping Close attached to
// the real body so the connection is still released.
type replayBody struct {
	io.Reader
	io.Closer
}

// WrapHTTPError creates a new HTTPError object from the given request/response
// ands wraps an existing error
func WrapHTTPError(err error, request *http.Request, response *http.Response) HTTPError {

	result := NewHTTPError(request, response)
	result.WrappedValue = err

	return result
}

// RedactedValue replaces a header value that carries a credential. It is exported so
// that a caller formatting its own output marks a redaction the same way derp does.
const RedactedValue = "[REDACTED]"

// redactedHeaders are the headers whose values are replaced when a transaction is
// recorded, keyed in canonical form.
//
// RULE: this map stays unexported, and callers reach it through IsSensitiveHeader.
// An exported map can be emptied by anyone who imports derp, which would silently
// disable redaction for the whole process.
var redactedHeaders = map[string]bool{
	"Authorization":       true,
	"Cookie":              true,
	"Proxy-Authorization": true,
	"Set-Cookie":          true,
	"Signature":           true,
	"X-Api-Key":           true,
}

// IsSensitiveHeader returns TRUE if a header carries a credential that must not be
// logged. Use it when formatting headers by hand; use RedactHeader when you have a map.
func IsSensitiveHeader(name string) bool {
	return redactedHeaders[textproto.CanonicalMIMEHeaderKey(name)]
}

// RedactHeader returns a copy of the provided headers with credential-bearing
// values replaced
func RedactHeader(header http.Header) http.Header {

	// RULE: copy before writing. These headers belong to a live http.Request or
	// http.Response, and redacting in place would strip the credential from the very
	// request that is still using it.
	result := header.Clone()

	if result == nil {
		return nil
	}

	// RULE: walk the map's OWN keys rather than looking up canonical names. A Header
	// built by hand can hold `authorization`, which no canonical lookup finds -- and
	// JSON serializes the raw map, so a missed key is a published credential.
	for name := range result {
		if redactedHeaders[textproto.CanonicalMIMEHeaderKey(name)] {
			result[name] = []string{RedactedValue}
		}
	}

	// A recorded error outlives the transaction it describes: it is wrapped, returned,
	// serialized, and written to a log somebody else reads.
	return result
}

// HTTPRequestReport includes details of a failed HTTP request
type HTTPRequestReport struct {
	URL    string      `json:"url"`    // Fully qualified URL that was requested
	Method string      `json:"method"` // HTTP method (GET, POST, etc) used to make the request
	Header http.Header `json:"header"` // Headers sent with the request, with credential-bearing values redacted
}

// HTTPResponseReport includes response details of a failed HTTP request
type HTTPResponseReport struct {
	StatusCode    int         `json:"statusCode"`              // Numeric HTTP status code returned by the server
	Status        string      `json:"status"`                  // Human-readable status line returned by the server
	Header        http.Header `json:"header"`                  // Headers returned with the response, with credential-bearing values redacted
	Body          string      `json:"body,omitempty"`          // Up to the first 4KB of the response body, which is usually where the server says what was wrong
	BodyTruncated bool        `json:"bodyTruncated,omitempty"` // TRUE when the body was longer than the cap, so a reader knows it continues
}

// Error implements the Error interface, which allows derp.Error objects to be
// used anywhere a standard error is used.
func (err HTTPError) Error() string {
	return err.Response.Status
}

// GetErrorCode returns the error Code embedded in this Error.
func (err HTTPError) GetErrorCode() int {
	return err.Response.StatusCode
}

// Unwrap returns the inner error wrapped by this HTTPError.
func (err HTTPError) Unwrap() error {
	return err.WrappedValue
}

/******************************************
 * HTTP-Specific Helper Methods
 *****************************************/

// GetRetryAfter returns the number of seconds to wait until retrying
// the transaction.  It is derived from one of several possible headers
// in the HTTP response, including `Retry-After`, `X-Ratelimit-Reset`,
// and `X-Rate-Limit-Reset`.
//
// If no such header is found, this method returns a default of 1 hour.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/429
func (err HTTPError) GetRetryAfter() time.Duration {

	// List of headers that might contain retry-after information
	headers := []string{
		"Retry-After",        // https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Retry-After
		"X-Ratelimit-Reset",  // https://www.ietf.org/archive/id/draft-polli-ratelimit-headers-02.html
		"X-Rate-Limit-Reset", // Some APIs use this variation
	}

	// Try each header in the list
	for _, header := range headers {

		// Get the header value
		value := err.Response.Header.Get(header)

		// If the header is empty, then skip
		if value == "" {
			continue
		}

		// Integers represent the number of seconds to wait
		// (named `atoiError` so that it does not shadow the `err` receiver)
		if seconds, atoiError := strconv.Atoi(value); atoiError == nil {
			return nonNegative(time.Duration(seconds) * time.Second)
		}

		// RFC3339 timestamps represent the time when the rate limit resets
		if resetAt, parseError := time.Parse(time.RFC3339, value); parseError == nil {
			return nonNegative(time.Until(resetAt))
		}

		// RFC1123 timestamps represent the time when the rate limit resets
		if resetAt, parseError := time.Parse(time.RFC1123, value); parseError == nil {
			return nonNegative(time.Until(resetAt))
		}

		// Last resort: the remaining legal HTTP-date formats (RFC850 and asctime), which
		// RFC9110 requires recipients to accept.
		// https://www.rfc-editor.org/rfc/rfc9110.html#name-date-time-formats
		if resetAt, parseError := http.ParseTime(value); parseError == nil {
			return nonNegative(time.Until(resetAt))
		}
	}

	// If no value is found, wait 1 hour before retrying
	return time.Hour
}

// nonNegative clamps a retry-after duration to zero.
func nonNegative(duration time.Duration) time.Duration {

	// A deadline that has already passed means the limit has reset, so wait no longer.
	if duration < 0 {
		return 0
	}

	return duration
}
