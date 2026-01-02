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

	require.Equal(t, time.Duration(120)*time.Second, httpError.GetRetryAfter().Truncate(time.Second))
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

	require.Equal(t, time.Duration(120)*time.Second, httpError.GetRetryAfter().Truncate(time.Second))
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

	require.Equal(t, time.Duration(120)*time.Second, httpError.GetRetryAfter().Truncate(time.Second))
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

	require.Equal(t, time.Duration(120)*time.Second, httpError.GetRetryAfter().Truncate(time.Second))
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

	require.Equal(t, time.Duration(120)*time.Second, httpError.GetRetryAfter().Truncate(time.Second))
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

	require.Equal(t, time.Duration(120)*time.Second, httpError.GetRetryAfter().Truncate(time.Second))
}

func TestHTTPRetry_RetryAfterWrapped(t *testing.T) {

	err := Error{
		WrappedValue: HTTPError{
			Response: HTTPResponseReport{
				StatusCode: 429,
				Header: http.Header{
					"Retry-After": []string{"120"},
				},
			},
		},
	}

	require.Equal(t, time.Duration(120)*time.Second, err.GetRetryAfter())
}

func TestHTTPRetry_RetryAfterWrappedGlobalFunc(t *testing.T) {

	err := Error{
		Code: 429,
		WrappedValue: HTTPError{
			Response: HTTPResponseReport{
				StatusCode: 429,
				Header: http.Header{
					"Retry-After": []string{"120"},
				},
			},
		},
	}

	ok, replyAfter := IsTooManyRequests(err)

	require.True(t, ok)
	require.Equal(t, time.Duration(120)*time.Second, replyAfter.Truncate(time.Second))
}
