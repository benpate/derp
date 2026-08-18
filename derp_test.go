package derp

import (
	"encoding/json"
	"errors"
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

func TestWrapGenericError(t *testing.T) {

	generic := errors.New("oof. that was bad")
	err := Wrap(generic, "TestEmptyWrappedValue", "Don't Do This").(Error)

	assert.Equal(t, 500, err.Code)
	assert.NotNil(t, err.WrappedValue)
	assert.Equal(t, "TestEmptyWrappedValue", err.Location)
	assert.Equal(t, "Don't Do This", err.Message)
	// assert.Equal(t, len(err.Details), 1)

	unwrapped := err.Unwrap()
	assert.Equal(t, "oof. that was bad", unwrapped.Error())
	Report(err)
}

func TestWrap_EmptyValue(t *testing.T) {

	{
		err := Wrap(nil, "TestEmptyWrappedValue", "This will still return an error")
		assert.Error(t, err)
	}

	{
		var innerError error
		outer := Wrap(innerError, "Should Still Return an error value", "really")
		assert.Error(t, outer)
	}
}

func TestWrapIF_EmptyValue(t *testing.T) {

	{
		err := WrapIF(nil, "TestEmptyWrappedValue", "This should return nil")
		assert.Nil(t, err)
	}

	{
		var innerError error
		outer := WrapIF(innerError, "Should Still Be Empty", "Really")
		assert.Nil(t, outer)
	}
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

func (w weirdErrorType) Error() string {
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

func (a arrayError) Error() string { return "array error" }

// funcError has Kind "Func", which CAN be nil.
type funcError func()

func (f funcError) Error() string { return "func error" }

// mapError has Kind "Map", which CAN be nil.
type mapError map[string]string

func (m mapError) Error() string { return "map error" }

// sliceError has Kind "Slice", which CAN be nil.
type sliceError []string

func (s sliceError) Error() string { return "slice error" }

// chanError has Kind "Chan", which CAN be nil.
type chanError chan int

func (c chanError) Error() string { return "chan error" }

// structError has Kind "Struct", which can NEVER be nil.
type structError struct{}

func (s structError) Error() string { return "struct error" }

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

// TestWrap_TypedNilPointer confirms that wrapping a typed-nil *Error does not
// call Error() on the nil pointer.
func TestWrap_TypedNilPointer(t *testing.T) {

	var inner *Error

	// Wrap always returns a non-nil error, even around a typed-nil inner value
	wrapped := Wrap(inner, "Location", "Message")
	require.NotNil(t, wrapped)
	require.Equal(t, "Location: Message", wrapped.Error())

	// The typed-nil inner value is not serialized into the Details
	require.Empty(t, Details(wrapped))

	// WrapIF treats a typed-nil inner value as nil, and returns nil
	require.Nil(t, WrapIF(inner, "Location", "Message"))
}

func TestNilErrorCode(t *testing.T) {
	require.Equal(t, 0, ErrorCode(nil))
}

func TestReportAndReturn(t *testing.T) {

	{
		err := errors.New("regular error")
		require.Equal(t, err, ReportAndReturn(err))
	}

	{
		err := newError(404, "Location", "Message")
		require.Equal(t, err, ReportAndReturn(err))
	}
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

func TestWrapIF_NotNil(t *testing.T) {
	err := WrapIF(errors.New("inner"), "location", "message")
	require.Error(t, err)
	require.Equal(t, codeInternalError, ErrorCode(err))
}

/******************************************
 * Fuzz Tests
 ******************************************/

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

// FuzzWrapChain confirms that arbitrarily deep chains of wrapped errors can be
// unwrapped without panicking, and that the root values come from the innermost error.
func FuzzWrapChain(f *testing.F) {

	f.Add(uint8(0), "root", "outer")
	f.Add(uint8(1), "", "")
	f.Add(uint8(10), "root message", "wrapper")
	f.Add(uint8(255), "deep", "d")

	f.Fuzz(func(t *testing.T, depth uint8, rootMessage string, wrapMessage string) {

		// Build a chain of `depth` wrappers around a single root error
		root := newError(codeNotFoundError, "RootLocation", rootMessage)

		var err error = root
		for index := 0; index < int(depth); index++ {
			err = Wrap(err, "WrapLocation", wrapMessage)
		}

		// Wrap NEVER returns nil
		require.NotNil(t, err)

		// The code is inherited all the way up the chain
		require.Equal(t, codeNotFoundError, ErrorCode(err))

		// Unwrapping always terminates at the root error
		require.Equal(t, root.Error(), Unwrap(err).Error())
		require.Equal(t, RootCause(err), Unwrap(err))

		// The root location is always the innermost location
		require.Equal(t, "RootLocation", RootLocation(err))

		// The root message is the innermost NON-EMPTY message, so an empty root
		// message falls back to the message of the shallowest wrapper that has one.
		if rootMessage != "" {
			require.Equal(t, rootMessage, RootMessage(err))
		}
	})
}
