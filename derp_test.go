package derp

import (
	"errors"
	"net/http"
	"testing"
	"time"

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
