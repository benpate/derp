package derp

import (
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestHTTPRetry_RetryAfterSeconds(t *testing.T) {

	httpError := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header: http.Header{
				"Retry-After": []string{"120"},
			},
		},
	}

	require.Equal(t, 120, httpError.RetryAfter())
}

func TestHTTPRetry_RateLimitSeconds(t *testing.T) {

	httpError := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header: http.Header{
				"X-Ratelimit-Reset": []string{"120"},
			},
		},
	}

	require.Equal(t, 120, httpError.RetryAfter())
}
func TestHTTPRetry_Rate_LimitSeconds(t *testing.T) {

	httpError := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header: http.Header{
				"X-Rate-Limit-Reset": []string{"120"},
			},
		},
	}

	require.Equal(t, 120, httpError.RetryAfter())
}

func TestHTTPRetry_RetryAfterTimestamp(t *testing.T) {

	timestamp := time.Now().Add(121 * time.Second).Format(time.RFC3339)

	httpError := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header: http.Header{
				"Retry-After": []string{timestamp},
			},
		},
	}

	require.Equal(t, 120, httpError.RetryAfter())
}

func TestHTTPRetry_RateLimitTimestamp(t *testing.T) {

	timestamp := time.Now().Add(121 * time.Second).Format(time.RFC3339)

	httpError := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header: http.Header{
				"X-Ratelimit-Reset": []string{timestamp},
			},
		},
	}

	require.Equal(t, 120, httpError.RetryAfter())
}

func TestHTTPRetry_Rate_LimitTimestamp(t *testing.T) {

	timestamp := time.Now().Add(121 * time.Second).Format(time.RFC3339)

	httpError := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header: http.Header{
				"X-Rate-Limit-Reset": []string{timestamp},
			},
		},
	}

	require.Equal(t, 120, httpError.RetryAfter())
}
