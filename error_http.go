package derp

import (
	"net/http"
	"strconv"
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

// RetryAfter returns the integer value of the `Retry-After` header,
// which is the number of seconds that the server should wait before
// retying a "Too Many Requests" response.
// If the `Retry-After` header is a valid integer, then this method
// returns the converted integer value. Otherwise, it returns zero.
// https://developer.mozilla.org/en-US/docs/Web/HTTP/Reference/Status/429
func (err HTTPError) RetryAfter() int {

	if retryAfter := err.Response.Header.Get("Retry-After"); retryAfter != "" {
		if retryAfterInt, err := strconv.Atoi(retryAfter); err == nil {
			return retryAfterInt
		}
	}
	return 0
}
