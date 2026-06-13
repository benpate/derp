package derp

import (
	"errors"
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

func TestNewHTTPError(t *testing.T) {

	// nil request and response produce an empty (but usable) HTTPError
	{
		err := NewHTTPError(nil, nil)
		require.Equal(t, "", err.Request.URL)
		require.Equal(t, 0, err.Response.StatusCode)
	}

	// a populated request and response are copied into the report
	{
		request, newErr := http.NewRequest(http.MethodGet, "https://example.com/path", nil)
		require.NoError(t, newErr)

		response := &http.Response{
			StatusCode: 404,
			Status:     "404 Not Found",
			Header:     http.Header{"Content-Type": []string{"text/plain"}},
		}

		err := NewHTTPError(request, response)

		require.Equal(t, "https://example.com/path", err.Request.URL)
		require.Equal(t, http.MethodGet, err.Request.Method)
		require.Equal(t, 404, err.Response.StatusCode)
		require.Equal(t, "404 Not Found", err.Response.Status)
		require.Equal(t, "text/plain", err.Response.Header.Get("Content-Type"))
	}
}

func TestWrapHTTPError(t *testing.T) {

	inner := errors.New("inner error")
	response := &http.Response{StatusCode: 500, Status: "500 Internal Server Error"}

	err := WrapHTTPError(inner, nil, response)

	require.Equal(t, inner, err.WrappedValue)
	require.Equal(t, inner, err.Unwrap())
}

func TestHTTPError_Interface(t *testing.T) {

	err := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 404,
			Status:     "404 Not Found",
		},
		WrappedValue: errors.New("inner error"),
	}

	require.Equal(t, "404 Not Found", err.Error())
	require.Equal(t, 404, err.GetErrorCode())
	require.Equal(t, "inner error", err.Unwrap().Error())
}

func TestHTTPError_GetRetryAfter_Default(t *testing.T) {

	// With no recognized headers, GetRetryAfter falls back to 1 hour.
	err := HTTPError{
		Response: HTTPResponseReport{StatusCode: 429},
	}
	require.Equal(t, time.Hour, err.GetRetryAfter())
}

func TestHTTPError_GetRetryAfter_InvalidValue(t *testing.T) {

	// A header value that is neither an integer nor a timestamp is ignored,
	// so GetRetryAfter falls back to 1 hour.
	err := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header:     http.Header{"Retry-After": []string{"not-a-number"}},
		},
	}
	require.Equal(t, time.Hour, err.GetRetryAfter())
}

func TestHTTPError_GetRetryAfter_RFC1123(t *testing.T) {

	timestamp := time.Now().Add(121 * time.Second).Format(time.RFC1123)

	err := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header:     http.Header{"Retry-After": []string{timestamp}},
		},
	}
	require.Equal(t, 120*time.Second, err.GetRetryAfter().Truncate(time.Second))
}
