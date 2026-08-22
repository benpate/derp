package derp

import (
	"errors"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

// Tests for the error-code predicates and nil-checking helpers in derp_codes.go

func TestIsNil(t *testing.T) {

	// IsNil has some strange edge cases, so make sure that nobody
	// makes derp panic because they define a strange error type

	var nilPointer *Error
	require.True(t, IsNil(nilPointer))

	var nilInterface error
	require.True(t, IsNil(nilInterface))

	actualError := errors.New("this should not be nil")
	require.False(t, IsNil(actualError))

	derpError := newError(404, "Code Location", "Error Message")
	require.False(t, IsNil(derpError))
}

func TestNotNil(t *testing.T) {

	var nilPointer *Error
	require.False(t, NotNil(nilPointer))

	var nilInterface error
	require.False(t, NotNil(nilInterface))

	actualError := errors.New("this should not be nil")
	require.True(t, NotNil(actualError))

	derpError := newError(0, "Code Location", "Error Message")
	require.True(t, NotNil(derpError))
}

func TestNotFoundOrGone(t *testing.T) {

	{
		require.False(t, IsNotFoundOrGone(nil))
	}

	{
		err := newError(500, "", "")
		require.False(t, IsNotFoundOrGone(err))
	}

	{
		err := newError(404, "", "")
		require.Equal(t, codeNotFoundError, ErrorCode(err))
		require.True(t, IsNotFoundOrGone(err))
	}

	{
		err := newError(410, "", "")
		require.Equal(t, codeGoneError, ErrorCode(err))
		require.True(t, IsNotFoundOrGone(err))
	}

	{
		err := errors.New("not found")
		require.True(t, IsNotFoundOrGone(err))
	}
}

type weirdErrorType string

func (weirdErrorType) Error() string {
	return "sure, it's an error"
}

func TestIsNil_WeirdErrorTypes(t *testing.T) {
	{
		require.False(t, IsNil(weirdErrorType("")))
	}
}

// The types below define errors whose underlying reflect.Kind varies.  IsNil
// must return a correct answer for every one of them, and must never panic --
// see https://medium.com/@mangatmodi/go-check-nil-interface-the-right-way-d142776edef1

// arrayError has Kind "Array", which can NEVER be nil.
type arrayError [2]byte

func (arrayError) Error() string { return "array error" }

// funcError has Kind "Func", which CAN be nil.
type funcError func()

func (funcError) Error() string { return "func error" }

// mapError has Kind "Map", which CAN be nil.
type mapError map[string]string

func (mapError) Error() string { return "map error" }

// sliceError has Kind "Slice", which CAN be nil.
type sliceError []string

func (sliceError) Error() string { return "slice error" }

// chanError has Kind "Chan", which CAN be nil.
type chanError chan int

func (chanError) Error() string { return "chan error" }

// structError has Kind "Struct", which can NEVER be nil.
type structError struct{}

func (structError) Error() string { return "struct error" }

// TestIsNil_AllKinds confirms that IsNil handles every Kind of error that a
// caller might define.  Kinds that cannot be nil (Array, Struct, string) must
// report NOT nil instead of panicking.
func TestIsNil_AllKinds(t *testing.T) {

	// Errors that are nil, and must be reported as nil
	nilCases := map[string]error{
		"nil func":    funcError(nil),
		"nil map":     mapError(nil),
		"nil slice":   sliceError(nil),
		"nil chan":    chanError(nil),
		"nil pointer": (*Error)(nil),
	}

	for name, err := range nilCases {
		t.Run(name, func(t *testing.T) {
			require.True(t, IsNil(err), "%s should be nil", name)
			require.False(t, NotNil(err), "%s should not be NotNil", name)
		})
	}

	// Errors that are NOT nil, and must be reported as not-nil (without panicking)
	notNilCases := map[string]error{
		"array":           arrayError{},
		"struct":          structError{},
		"string":          weirdErrorType(""),
		"populated func":  funcError(func() {}),
		"populated map":   mapError{},
		"populated slice": sliceError{},
		"populated chan":  chanError(make(chan int)),
		"populated ptr":   &Error{},
	}

	for name, err := range notNilCases {
		t.Run(name, func(t *testing.T) {
			require.NotPanics(t, func() { IsNil(err) }, "%s must not panic", name)
			require.False(t, IsNil(err), "%s should not be nil", name)
			require.True(t, NotNil(err), "%s should be NotNil", name)
		})
	}
}

// TestPublicAPI_ArrayError confirms that an Array-Kinded error (which cannot be
// nil, and once panicked inside IsNil) travels safely through the public API.
func TestPublicAPI_ArrayError(t *testing.T) {

	err := arrayError{}

	require.NotPanics(t, func() {
		require.Equal(t, "array error", Message(err))
		require.Equal(t, codeInternalError, ErrorCode(err))
		require.Equal(t, "", Location(err))
		require.Nil(t, Details(err))
		require.Zero(t, RetryAfter(err))
		require.Equal(t, "array error", RootMessage(err))
		Report(err)
	})
}

func TestIsInformational(t *testing.T) {
	{
		e := newError(99, "location", "message")
		require.False(t, IsInformational(e))
	}
	{
		e := newError(100, "Location", "Message")
		require.True(t, IsInformational(e))
	}
	{
		e := newError(199, "Location", "Message")
		require.True(t, IsInformational(e))
	}
	{
		e := newError(200, "Location", "Message")
		require.False(t, IsInformational(e))
	}
}

func TestIsSuccess(t *testing.T) {
	{
		e := newError(199, "location", "message")
		require.False(t, IsSuccess(e))
	}
	{
		e := newError(200, "Location", "Message")
		require.True(t, IsSuccess(e))
	}
	{
		e := newError(299, "Location", "Message")
		require.True(t, IsSuccess(e))
	}
	{
		e := newError(300, "Location", "Message")
		require.False(t, IsSuccess(e))
	}
}

func TestIsRedirection(t *testing.T) {
	{
		e := newError(299, "location", "message")
		require.False(t, IsRedirection(e))
	}
	{
		e := newError(300, "Location", "Message")
		require.True(t, IsRedirection(e))
	}
	{
		e := newError(399, "Location", "Message")
		require.True(t, IsRedirection(e))
	}
	{
		e := newError(400, "Location", "Message")
		require.False(t, IsRedirection(e))
	}
}

func TestIsClientError(t *testing.T) {
	{
		e := newError(399, "location", "message")
		require.False(t, IsClientError(e))
	}
	{
		e := newError(400, "Location", "Message")
		require.True(t, IsClientError(e))
	}
	{
		e := newError(499, "Location", "Message")
		require.True(t, IsClientError(e))
	}
	{
		e := newError(500, "Location", "Message")
		require.False(t, IsClientError(e))
	}
}

func TestIsServerError(t *testing.T) {
	{
		e := newError(499, "location", "message")
		require.False(t, IsServerError(e))
	}
	{
		e := newError(500, "Location", "Message")
		require.True(t, IsServerError(e))
	}
	{
		e := newError(599, "Location", "Message")
		require.True(t, IsServerError(e))
	}
	{
		e := newError(600, "Location", "Message")
		require.False(t, IsServerError(e))
	}
}

func TestIsBadRequest(t *testing.T) {

	otherError := newError(0, "location", "message")
	require.False(t, IsBadRequest(otherError))

	badRequest := newError(400, "Location", "Message")
	require.True(t, IsBadRequest(badRequest))
}

func TestIsUnauthorized(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsUnauthorized(otherError))

	unauthorized := newError(401, "Location", "Message")
	require.True(t, IsUnauthorized(unauthorized))
}

func TestIsForbidden(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsForbidden(otherError))

	forbidden := newError(403, "Location", "Message")
	require.True(t, IsForbidden(forbidden))
}

func TestIsNotFound(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsNotFound(otherError))

	notFoundCode := newError(404, "Location", "Message")
	require.True(t, IsNotFound(notFoundCode))

	notFoundText := errors.New("not found")
	require.True(t, IsNotFound(notFoundText))
}

func TestIsConflict(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsConflict(otherError))

	conflict := Conflict("Location", "Message")
	require.True(t, IsConflict(conflict))

	// The code must survive Wrap, so a data-layer conflict stays a conflict up the stack.
	require.True(t, IsConflict(Wrap(conflict, "outer", "wrapped")))

	// WithConflict sets the code on any constructor.
	require.True(t, IsConflict(newError(0, "location", "message", WithConflict())))
}

func TestIsTeapot(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsTeapot(otherError))

	teapot := newError(418, "Location", "Message")
	require.True(t, IsTeapot(teapot))
}

func TestIsMisdirectedRequest(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsMisdirectedRequest(otherError))

	misdirected := newError(421, "Location", "Message")
	require.True(t, IsMisdirectedRequest(misdirected))
}

func TestIsValidationError(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsValidationError(otherError))

	validation := newError(422, "Location", "Message")
	require.True(t, IsValidationError(validation))
}

func TestIsInternalServerError(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsInternalServerError(otherError))

	internal := newError(500, "Location", "Message")
	require.True(t, IsInternalServerError(internal))
}

func TestIsNotImplemented(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsNotImplemented(otherError))

	notImplemented := newError(501, "Location", "Message")
	require.True(t, IsNotImplemented(notImplemented))
}

func TestIsBadGateway(t *testing.T) {
	otherError := newError(0, "location", "message")
	require.False(t, IsBadGateway(otherError))

	badGateway := newError(502, "Location", "Message")
	require.True(t, IsBadGateway(badGateway))
}

func TestIsGone(t *testing.T) {
	require.False(t, IsGone(newError(500, "location", "message")))
	require.True(t, IsGone(newError(codeGoneError, "location", "message")))
}

func TestIsTooManyRequests(t *testing.T) {

	// a non-429 error is not "too many requests"
	{
		ok, retryAfter := IsTooManyRequests(newError(404, "location", "message"))
		require.False(t, ok)
		require.Equal(t, time.Duration(0), retryAfter)
	}

	// a 429 error with no retry-after header defaults to 1 hour
	{
		ok, retryAfter := IsTooManyRequests(newError(codeTooManyRequestsError, "location", "message"))
		require.True(t, ok)
		require.Equal(t, time.Hour, retryAfter)
	}

	// a 429 error with a retry-after header uses that duration
	{
		err := Error{
			Code: codeTooManyRequestsError,
			WrappedValue: HTTPError{
				Response: HTTPResponseReport{
					Header: http.Header{"Retry-After": []string{"120"}},
				},
			},
		}
		ok, retryAfter := IsTooManyRequests(err)
		require.True(t, ok)
		require.Equal(t, 120*time.Second, retryAfter.Truncate(time.Second))
	}
}

// FuzzErrorCodeRanges confirms that the range predicates partition the space of
// error codes: no code may ever belong to two different ranges.
func FuzzErrorCodeRanges(f *testing.F) {

	f.Add(0)
	f.Add(100)
	f.Add(199)
	f.Add(200)
	f.Add(299)
	f.Add(399)
	f.Add(400)
	f.Add(499)
	f.Add(500)
	f.Add(599)
	f.Add(600)
	f.Add(-1)
	f.Add(math.MaxInt32)
	f.Add(math.MinInt32)

	f.Fuzz(func(t *testing.T, code int) {

		// Use an empty message so that the "not found" text fallback cannot apply
		err := newError(code, "", "")

		// Count how many ranges claim this code.  The answer is never more than one.
		ranges := []bool{
			IsInformational(err),
			IsSuccess(err),
			IsRedirection(err),
			IsClientError(err),
			IsServerError(err),
		}

		matches := 0
		for _, isMatch := range ranges {
			if isMatch {
				matches++
			}
		}

		require.LessOrEqual(t, matches, 1, "code %d matched %d ranges", code, matches)

		// Codes outside 100-599 belong to no range at all
		if (code < 100) || (code >= 600) {
			require.Zero(t, matches, "code %d should match no range", code)
		}

		// Specific codes always imply their enclosing range
		if IsNotFound(err) || IsGone(err) || IsConflict(err) || IsTeapot(err) {
			require.True(t, IsClientError(err), "code %d should be a client error", code)
		}

		if IsInternalServerError(err) || IsNotImplemented(err) || IsBadGateway(err) {
			require.True(t, IsServerError(err), "code %d should be a server error", code)
		}
	})
}

// FuzzIsNotFound confirms that the "not found" message fallback tolerates any
// message text, and that IsNotFound always implies IsNotFoundOrGone.
func FuzzIsNotFound(f *testing.F) {

	f.Add("not found")
	f.Add("NOT FOUND")
	f.Add("Not Found")
	f.Add("")
	f.Add("mongo: no documents in result")
	f.Add("\xff\xfe not found")
	f.Add(strings.Repeat("not found", 128))

	f.Fuzz(func(t *testing.T, message string) {

		err := errors.New(message)

		// IsNotFound always implies IsNotFoundOrGone
		if IsNotFound(err) {
			require.True(t, IsNotFoundOrGone(err), "message %q", message)
		}

		// The fallback is a case-insensitive match on exactly "not found"
		require.Equal(t, strings.EqualFold(message, "not found"), IsNotFound(err), "message %q", message)
	})
}
