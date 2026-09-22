package derp

import (
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"net/http"
	"strings"
	"testing"
	"time"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestDerp tests basic derp functions (separate from features of a specific reporter)
func TestDerp(t *testing.T) {

	// Create an inner error
	innerError := NotFound("WrappedValue", "Not Found", "detail1", "detail2", "detail3")

	assert.Equal(t, innerError.Location, "WrappedValue")
	assert.Equal(t, innerError.Message, "Not Found")
	assert.Equal(t, innerError.Code, 404)
	assert.Equal(t, innerError.Details[0], "detail1")
	assert.Equal(t, innerError.Details[1], "detail2")
	assert.Equal(t, innerError.Details[2], "detail3")
	assert.Equal(t, IsNotFound(innerError), true)

	// Create an outer error
	outerError := Wrap(innerError, "OuterError", "Inherited", "other details here").(Error)

	assert.Equal(t, outerError.Location, "OuterError")
	assert.Equal(t, outerError.Message, "Inherited")
	assert.Equal(t, outerError.Code, 404) // This is still 404 because we've let the inner error code bubble up
	assert.NotNil(t, outerError.WrappedValue)
	assert.Equal(t, outerError.Details[0], "other details here")
	assert.Equal(t, IsNotFound(outerError), true)

	// Test the RootCause() function
	assert.Equal(t, "WrappedValue", RootCause(outerError).(Error).Location)
}

func TestNewConvenienceFns(t *testing.T) {

	badRequest := BadRequest("location", "description")
	require.Equal(t, codeBadRequestError, ErrorCode(badRequest))

	forbidden := Forbidden("location", "description")
	require.Equal(t, codeForbiddenError, ErrorCode(forbidden))

	internal := Internal("location", "description")
	require.Equal(t, codeInternalError, ErrorCode(internal))

	notFound := NotFound("location", "description")
	require.Equal(t, codeNotFoundError, ErrorCode(notFound))

	unauthorized := Unauthorized("location", "description")
	require.Equal(t, codeUnauthorizedError, ErrorCode(unauthorized))

	invalid := Validation("location", "description")
	require.Equal(t, codeValidationError, ErrorCode(invalid))

	teapot := Teapot("location", "description")
	require.Equal(t, codeTeapotError, ErrorCode(teapot))

	misdirected := MisdirectedRequest("location", "description")
	require.Equal(t, codeMisdirectedRequestError, ErrorCode(misdirected))

	notImplemented := NotImplemented("location", "description")
	require.Equal(t, codeNotImplementedError, ErrorCode(notImplemented))

	badGateway := BadGateway("location", "description")
	require.Equal(t, codeBadGatewayError, ErrorCode(badGateway))
}

func TestMessage(t *testing.T) {

	require.Equal(t, "", Message(nil))

	derp := NotFound("location", "description")
	require.Equal(t, "description", Message(derp))

	standard := errors.New("this is a standard error")
	require.Equal(t, "this is a standard error", Message(standard))
}

func TestErrorInterface(t *testing.T) {

	// Create an error
	innerError := NotFound("Location Name", "Error Description", "details")

	// Verify that the error interface is outputting what we expect.
	assert.Equal(t, innerError.Error(), "Location Name: Error Description")
}

func TestStandardError(t *testing.T) {

	// Testing how derp handles an error from the standard library
	err := errors.New("This is a standard error")

	// Wrap it the stdlib error in a derp.  This means: 1) assigning an error code, and 2) making the original error message a property of the derp.Error.
	outer := Internal("TestStandardError", "Encapsulating Error", err.Error())

	assert.Equal(t, "TestStandardError", outer.Location)
	assert.Equal(t, "Encapsulating Error", outer.Message)
	assert.Equal(t, 500, outer.Code)
	assert.Equal(t, 1, len(outer.Details))
	assert.Nil(t, outer.WrappedValue)
}

func TestNotFound(t *testing.T) {

	require.False(t, IsNotFound(nil))

	{
		err := errors.New("regular error")
		require.False(t, IsNotFound(err))
	}

	{
		err := errors.New("not found")
		require.True(t, IsNotFound(err))
	}

	{
		err := newError(500, "", "")
		require.False(t, IsNotFound(err))
	}

	{
		err := newError(404, "", "")
		require.True(t, IsNotFound(err))
	}

	{
		e := NotFound("Location", "Message")
		assert.Equal(t, 404, ErrorCode(e))
	}
}

func TestNilErrorCode(t *testing.T) {
	require.Equal(t, 0, ErrorCode(nil))
}

// TestBadGateway confirms that a 502 is a server error, but is distinct from
// a generic 500 -- the whole point of adding it.
func TestBadGateway(t *testing.T) {
	err := BadGateway("location", "message")
	require.Equal(t, codeBadGatewayError, ErrorCode(err))
	require.True(t, IsBadGateway(err))
	require.True(t, IsServerError(err))
	require.False(t, IsInternalServerError(err))
}

func TestGone(t *testing.T) {
	err := Gone("location", "message")
	require.Equal(t, codeGoneError, ErrorCode(err))
	require.True(t, IsGone(err))
}

func TestTimeout(t *testing.T) {
	err := Timeout("location", "message")
	require.Equal(t, codeTimeout, ErrorCode(err))
}

// TestDeprecatedErrorFns exercises the deprecated *Error() constructors,
// confirming they still produce the same error codes as their replacements.
func TestDeprecatedErrorFns(t *testing.T) {
	require.Equal(t, codeBadRequestError, ErrorCode(BadRequestError("location", "message")))
	require.Equal(t, codeUnauthorizedError, ErrorCode(UnauthorizedError("location", "message")))
	require.Equal(t, codeForbiddenError, ErrorCode(ForbiddenError("location", "message")))
	require.Equal(t, codeMisdirectedRequestError, ErrorCode(MisdirectedRequestError("location", "message")))
	require.Equal(t, codeNotFoundError, ErrorCode(NotFoundError("location", "message")))
	require.Equal(t, codeTeapotError, ErrorCode(TeapotError("location", "message")))
	require.Equal(t, codeTimeout, ErrorCode(TimeoutError("location", "message")))
	require.Equal(t, codeValidationError, ErrorCode(ValidationError("message")))
	require.Equal(t, codeInternalError, ErrorCode(InternalError("location", "message")))
	require.Equal(t, codeNotImplementedError, ErrorCode(NotImplementedError("location")))
}

func TestLocation(t *testing.T) {

	// nil error has no location
	require.Equal(t, "", Location(nil))

	// standard errors do not implement LocationGetter
	require.Equal(t, "", Location(errors.New("standard error")))

	// derp errors return their location
	require.Equal(t, "the location", Location(newError(500, "the location", "message")))
}

func TestURL(t *testing.T) {

	// nil error has no URL
	require.Equal(t, "", URL(nil))

	// standard errors do not implement URLGetter
	require.Equal(t, "", URL(errors.New("standard error")))

	// derp errors return their URL
	require.Equal(t, "https://example.com/help", URL(Error{URL: "https://example.com/help"}))
}

func TestDetails(t *testing.T) {

	// nil error has no details
	require.Nil(t, Details(nil))

	// standard errors do not implement DetailsGetter
	require.Nil(t, Details(errors.New("standard error")))

	// derp errors return their details
	require.Equal(t, []any{"a", "b"}, Details(newError(500, "location", "message", "a", "b")))
}

func TestRetryAfter(t *testing.T) {

	// nil error has no retry-after
	require.Equal(t, time.Duration(0), RetryAfter(nil))

	// standard errors do not implement RetryAfterGetter
	require.Equal(t, time.Duration(0), RetryAfter(errors.New("standard error")))

	// HTTPError implements RetryAfterGetter
	httpError := HTTPError{
		Response: HTTPResponseReport{
			Header: http.Header{"Retry-After": []string{"120"}},
		},
	}
	require.Equal(t, 120*time.Second, RetryAfter(httpError).Truncate(time.Second))
}

func TestSerialize(t *testing.T) {

	// nil error serializes to an empty string
	require.Equal(t, "", Serialize(nil))

	// a valid error serializes to JSON
	require.Contains(t, Serialize(newError(404, "location", "message")), `"message":"message"`)

	// an error that cannot be marshaled (a func in Details) serializes to an empty string
	require.Equal(t, "", Serialize(Error{Details: []any{func() {}}}))
}

func TestRootMessage(t *testing.T) {

	// nil error has no message
	require.Equal(t, "", RootMessage(nil))

	// standard (non-unwrappable) errors return their own message
	require.Equal(t, "standard error", RootMessage(errors.New("standard error")))

	// a lone derp error returns its own message
	require.Equal(t, "only", RootMessage(newError(500, "location", "only")))

	// a chain returns the deepest non-empty message
	inner := newError(500, "InnerLocation", "InnerMessage")
	middle := Wrap(inner, "MiddleLocation", "MiddleMessage")
	outer := Wrap(middle, "OuterLocation", "OuterMessage")
	require.Equal(t, "InnerMessage", RootMessage(outer))

	// when deeper messages are empty, it falls back to the current message
	emptyInner := newError(500, "InnerLocation", "")
	wrapped := Wrap(emptyInner, "OuterLocation", "OuterMessage")
	require.Equal(t, "OuterMessage", RootMessage(wrapped))
}

func TestRootLocation(t *testing.T) {

	// nil error has no location
	require.Equal(t, "", RootLocation(nil))

	// standard (non-unwrappable) errors have no location
	require.Equal(t, "", RootLocation(errors.New("standard error")))

	// a lone derp error returns its own location
	require.Equal(t, "only", RootLocation(newError(500, "only", "message")))

	// a chain returns the deepest non-empty location
	inner := newError(500, "InnerLocation", "InnerMessage")
	middle := Wrap(inner, "MiddleLocation", "MiddleMessage")
	outer := Wrap(middle, "OuterLocation", "OuterMessage")
	require.Equal(t, "InnerLocation", RootLocation(outer))

	// when deeper locations are empty, it falls back to the current location
	emptyInner := newError(500, "", "InnerMessage")
	wrapped := Wrap(emptyInner, "OuterLocation", "OuterMessage")
	require.Equal(t, "OuterLocation", RootLocation(wrapped))
}

// FuzzSerialize confirms that any error survives a Serialize/Deserialize
// round trip without panicking, and without losing its numeric code.
func FuzzSerialize(f *testing.F) {

	f.Add(404, "location", "message", "detail")
	f.Add(0, "", "", "")
	f.Add(500, "derp.Test", "not found", "\x00\xff")
	f.Add(-1, "\xed\xa0\x80", "\u00e9\u00e8", "\U0001f600")
	f.Add(math.MaxInt32, "very long", strings.Repeat("x", 1024), "")

	f.Fuzz(func(t *testing.T, code int, location string, message string, detail string) {

		err := newError(code, location, message, detail)

		// Accessors always agree with the values used to build the error
		require.Equal(t, code, ErrorCode(err))
		require.Equal(t, location, Location(err))
		require.Equal(t, message, Message(err))
		require.Equal(t, []any{detail}, Details(err))

		// Serialize always produces parseable JSON
		serialized := Serialize(err)
		require.NotEmpty(t, serialized, "serializing %d / %q / %q", code, location, message)

		var parsed Error
		require.NoError(t, json.Unmarshal([]byte(serialized), &parsed))

		// The numeric code always round-trips exactly
		require.Equal(t, code, parsed.Code)

		// RULE: JSON replaces invalid UTF-8 with U+FFFD, so strings only
		// round-trip exactly when they were valid UTF-8 to begin with.
		if utf8.ValidString(location) && utf8.ValidString(message) {
			require.Equal(t, location, parsed.Location)
			require.Equal(t, message, parsed.Message)
		}
	})
}

func TestCodes(t *testing.T) {
	err := newError(123, "whatever", "dude")
	assert.Equal(t, 123, ErrorCode(err))
}

func TestCodeGeneric(t *testing.T) {
	err := errors.New("whatever, dude")
	assert.Equal(t, 500, ErrorCode(err))
}

// A foreign error type that wraps another error, standing in for text/template's ExecError.
type foreignWrapper struct {
	inner error
}

func (e foreignWrapper) Error() string { return "foreign: " + e.inner.Error() }
func (e foreignWrapper) Unwrap() error { return e.inner }

func TestErrorCode_SurvivesAForeignWrapper(t *testing.T) {

	// A derp error reached through a foreign wrapper keeps its own code.
	// Regression: html/template reports a failing method as its own ExecError, which
	// turned every 401 raised inside a template into a 500.
	unauthorized := Unauthorized("test.Location", "not signed in")

	require.Equal(t, 401, ErrorCode(unauthorized))
	require.Equal(t, 401, ErrorCode(foreignWrapper{inner: unauthorized}))
	require.Equal(t, 401, ErrorCode(fmt.Errorf("calling method: %w", unauthorized)))

	// Two foreign layers, which is what html/template actually produces
	// (ExecError wrapping a fmt.wrapError wrapping the original).
	doubled := foreignWrapper{inner: fmt.Errorf("calling method: %w", unauthorized)}
	require.Equal(t, 401, ErrorCode(doubled))

	// Wrapping through derp still works, and still reports the inner code.
	require.Equal(t, 401, ErrorCode(Wrap(doubled, "outer.Location", "executing template")))

	// An error with no code anywhere in the chain is still a 500.
	require.Equal(t, 500, ErrorCode(foreignWrapper{inner: errors.New("plain")}))
	require.Equal(t, 500, ErrorCode(errors.New("plain")))

	// Nil is still zero.
	require.Equal(t, 0, ErrorCode(nil))
}

// TestRetryAfter_SurvivesAForeignWrapper confirms that a 429's Retry-After is readable
// through any wrapper, not only through derp's own.
func TestRetryAfter_SurvivesAForeignWrapper(t *testing.T) {

	// RetryAfter once used a bare type assertion, so a non-derp layer zeroed the duration
	// while ErrorCode still read 429.  IsTooManyRequests then substituted its 1 hour
	// default, and a host asking for 30 seconds was deferred for an hour.
	rateLimited := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header:     http.Header{"Retry-After": []string{"30"}},
		},
	}

	require.Equal(t, 30*time.Second, RetryAfter(rateLimited))
	require.Equal(t, 30*time.Second, RetryAfter(foreignWrapper{inner: rateLimited}))
	require.Equal(t, 30*time.Second, RetryAfter(fmt.Errorf("transport: %w", rateLimited)))

	// Two foreign layers, and a foreign layer sandwiched inside derp wrapping, which is
	// the shape a queue consumer actually sees.
	doubled := foreignWrapper{inner: fmt.Errorf("transport: %w", rateLimited)}
	require.Equal(t, 30*time.Second, RetryAfter(doubled))
	require.Equal(t, 30*time.Second, RetryAfter(Wrap(doubled, "outer.Location", "loading")))
	require.Equal(t, 30*time.Second, RetryAfter(Wrap(Wrap(doubled, "a", "b"), "c", "d")))

	// errors.Join reaches the 429 through Unwrap() []error, which errors.As also walks.
	require.Equal(t, 30*time.Second, RetryAfter(errors.Join(errors.New("unrelated"), rateLimited)))

	// Nothing in the chain carries a duration, so there is none to report.
	require.Equal(t, time.Duration(0), RetryAfter(foreignWrapper{inner: errors.New("plain")}))
	require.Equal(t, time.Duration(0), RetryAfter(Wrap(nil, "location", "message")))
	require.Equal(t, time.Duration(0), RetryAfter(NotFound("location", "gone")))
	require.Equal(t, time.Duration(0), RetryAfter(nil))
}

// TestIsTooManyRequests_ReportsTheHostsOwnDelayThroughAnyWrapper confirms that the delay a
// caller reschedules on is the host's own, whatever wraps the error.
func TestIsTooManyRequests_ReportsTheHostsOwnDelayThroughAnyWrapper(t *testing.T) {

	// The caller-visible half of TestRetryAfter_SurvivesAForeignWrapper: a queue consumer
	// reschedules on this duration, so losing it changes behavior, not just log text.
	rateLimited := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header:     http.Header{"Retry-After": []string{"30"}},
		},
	}

	chains := map[string]error{
		"bare":            rateLimited,
		"derp.Wrap":       Wrap(rateLimited, "location", "message"),
		"foreign":         foreignWrapper{inner: rateLimited},
		"fmt.Errorf":      fmt.Errorf("transport: %w", rateLimited),
		"derp(foreign())": Wrap(foreignWrapper{inner: rateLimited}, "location", "message"),
	}

	for name, err := range chains {
		isTooMany, retryAfter := IsTooManyRequests(err)
		require.True(t, isTooMany, name)
		require.Equal(t, 30*time.Second, retryAfter, name)
	}

	// A 429 that sends no usable Retry-After still falls back to one hour, never zero.
	noHeader := HTTPError{Response: HTTPResponseReport{StatusCode: 429}}

	for name, err := range map[string]error{
		"bare":       noHeader,
		"fmt.Errorf": fmt.Errorf("transport: %w", noHeader),
	} {
		isTooMany, retryAfter := IsTooManyRequests(err)
		require.True(t, isTooMany, name)
		require.Equal(t, time.Hour, retryAfter, name)
	}
}

// TestErrorGetRetryAfter_StaysADelegation confirms that an Error reports the duration carried
// below it, which is what lets RetryAfter search past the outermost Error.
func TestErrorGetRetryAfter_StaysADelegation(t *testing.T) {

	// Pinned directly so that replacing the delegation with a field read fails HERE and
	// names the cause, rather than as several unrelated wrapped-chain assertions.
	inner := HTTPError{
		Response: HTTPResponseReport{
			StatusCode: 429,
			Header:     http.Header{"Retry-After": []string{"45"}},
		},
	}

	outer := Error{Code: 429, Location: "location", Message: "message", WrappedValue: inner}

	require.Equal(t, 45*time.Second, outer.GetRetryAfter())

	// An Error with nothing below it has no duration of its own to report.
	require.Equal(t, time.Duration(0), Error{Code: 429}.GetRetryAfter())
}
