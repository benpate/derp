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

// TestNewHTTPError_NoURL confirms that a Request with no URL is reported with
// an empty URL string, rather than dereferencing the nil URL.
func TestNewHTTPError_NoURL(t *testing.T) {

	request := &http.Request{Method: http.MethodGet}

	var result HTTPError
	require.NotPanics(t, func() {
		result = NewHTTPError(request, nil)
	})

	require.Equal(t, "", result.Request.URL)
	require.Equal(t, http.MethodGet, result.Request.Method)
}

// TestNewHTTPError_NilBoth confirms that a nil request AND a nil response
// produce an empty (but usable) HTTPError.
func TestNewHTTPError_NilBoth(t *testing.T) {

	result := NewHTTPError(nil, nil)

	require.Equal(t, "", result.Request.URL)
	require.Equal(t, "", result.Response.Status)
	require.Equal(t, 0, result.GetErrorCode())
	require.Equal(t, time.Hour, result.GetRetryAfter())
	require.Nil(t, result.Unwrap())
}

// FuzzGetRetryAfter feeds arbitrary, untrusted Retry-After header values to the parser.
// The contract under fuzzing is that parsing must never panic, and must never report a
// negative wait, no matter how malformed the input.
func FuzzGetRetryAfter(f *testing.F) {

	// Seed the corpus with values that exercise each parsing branch.
	f.Add("120")                           // integer seconds
	f.Add(time.Now().Format(time.RFC3339)) // RFC3339 timestamp
	f.Add(time.Now().Format(time.RFC1123)) // RFC1123 timestamp
	f.Add("not-a-number")                  // unparseable
	f.Add("")                              // empty
	f.Add("-1")                            // negative integer

	f.Fuzz(func(t *testing.T, value string) {
		httpError := HTTPError{
			Response: HTTPResponseReport{
				Header: http.Header{"Retry-After": []string{value}},
			},
		}

		// RULE: A retry-after duration is a wait, and a wait is never negative.
		// (a panic inside GetRetryAfter also fails the fuzz test automatically)
		if result := httpError.GetRetryAfter(); result < 0 {
			t.Fatalf("negative retry-after %v from header value %q", result, value)
		}
	})
}

// TestHTTPError_GetRetryAfter_NegativeSeconds confirms that a negative delay-seconds
// value is clamped to zero, instead of reporting a wait that runs backwards.
func TestHTTPError_GetRetryAfter_NegativeSeconds(t *testing.T) {

	err := HTTPError{
		Response: HTTPResponseReport{
			Header: http.Header{"Retry-After": []string{"-30"}},
		},
	}

	require.Equal(t, time.Duration(0), err.GetRetryAfter())
}

// TestHTTPError_GetRetryAfter_ExpiredTimestamp confirms that a reset time in the past
// is clamped to zero, because the rate limit has already reset.
func TestHTTPError_GetRetryAfter_ExpiredTimestamp(t *testing.T) {

	timestamp := time.Now().Add(-1 * time.Hour).Format(time.RFC3339)

	err := HTTPError{
		Response: HTTPResponseReport{
			Header: http.Header{"Retry-After": []string{timestamp}},
		},
	}

	require.Equal(t, time.Duration(0), err.GetRetryAfter())
}

// TestHTTPError_GetRetryAfter_RFC850 confirms that RFC850 dates -- one of the three
// HTTP-date formats that RFC9110 requires recipients to accept -- are parsed.
func TestHTTPError_GetRetryAfter_RFC850(t *testing.T) {

	timestamp := time.Now().UTC().Add(121 * time.Second).Format(time.RFC850)

	err := HTTPError{
		Response: HTTPResponseReport{
			Header: http.Header{"Retry-After": []string{timestamp}},
		},
	}

	require.Equal(t, 120*time.Second, err.GetRetryAfter().Truncate(time.Second))
}

// TestHTTPError_GetRetryAfter_ANSIC confirms that asctime dates -- one of the three
// HTTP-date formats that RFC9110 requires recipients to accept -- are parsed.
func TestHTTPError_GetRetryAfter_ANSIC(t *testing.T) {

	timestamp := time.Now().UTC().Add(121 * time.Second).Format(time.ANSIC)

	err := HTTPError{
		Response: HTTPResponseReport{
			Header: http.Header{"Retry-After": []string{timestamp}},
		},
	}

	require.Equal(t, 120*time.Second, err.GetRetryAfter().Truncate(time.Second))
}
