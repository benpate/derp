package derp

import (
	"net/http"
	"strconv"
	"time"
)

// HTTPError wraps a standard derp.Error, including additional data about a failed HTTP transaction
type HTTPError struct {
	Request  HTTPRequestReport  `json:"request"`
	Response HTTPResponseReport `json:"response"`

	WrappedValue error `json:"innerError,omitempty"` // An underlying error object used to identify the root cause of this error.
}

// NewHTTPError creates a new HTTPError object from the given request and response
func NewHTTPError(request *http.Request, response *http.Response) HTTPError {

	result := HTTPError{}

	if request != nil {
		result.Request = HTTPRequestReport{
			URL:    request.URL.String(),
			Method: request.Method,
			Header: request.Header,
		}
	}

	if response != nil {
		result.Response = HTTPResponseReport{
			StatusCode: response.StatusCode,
			Status:     response.Status,
			Header:     response.Header,
		}
	}

	return result
}

// WrapHTTPError creates a new HTTPError object from the given request/response
// ands wraps an existing error
func WrapHTTPError(err error, request *http.Request, response *http.Response) HTTPError {

	result := NewHTTPError(request, response)
	result.WrappedValue = err

	return result
}

// HTTPRequestReport includes details of a failed HTTP request
type HTTPRequestReport struct {
	URL    string      `json:"url"`
	Method string      `json:"method"`
	Header http.Header `json:"header"`
}

// HTTPResponseReport includes response details of a failed HTTP request
type HTTPResponseReport struct {
	StatusCode int         `json:"statusCode"`
	Status     string      `json:"status"`
	Header     http.Header `json:"header"`
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

// RetryAfter returns the number of seconds to wait until retrying
// the transaction.  It is derived from one of several possible headers
// in the HTTP response, including `Retry-After`, `X-Ratelimit-Reset`,
// and `X-Rate-Limit-Reset`.
//
// If no such header is found, this method returns a default of 3600 seconds (1 hour).
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/429
func (err HTTPError) RetryAfter() int {

	// List of headers that might contain retry-after information
	headers := []string{
		"Retry-After",       // https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Headers/Retry-After
		"X-Ratelimit-Reset", // https://www.ietf.org/archive/id/draft-polli-ratelimit-headers-02.html
		"X-Rate-Limit-Reset",
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
		if asInteger, err := strconv.Atoi(value); err == nil {
			return asInteger
		}

		// RFC3339 timestamps represent the time when the rate limit resets
		if asTimestamp, err := time.Parse(time.RFC3339, value); err == nil {
			return int(time.Until(asTimestamp).Seconds())
		}

		// RFC1123 timestamps represent the time when the rate limit resets
		if asTimestamp, err := time.Parse(time.RFC1123, value); err == nil {
			return int(time.Until(asTimestamp).Seconds())
		}
	}

	// If no value is found, wait 1 hour before retrying
	return int(time.Duration(time.Hour).Seconds())
}
