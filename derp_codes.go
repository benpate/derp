package derp

// NotFound returns TRUE if the error `Code` is a 404 / Not Found error.
func NotFound(err error) bool {

	if isNil(err) {
		return false
	}

	if coder, ok := err.(ErrorCodeGetter); ok {
		return coder.GetErrorCode() == CodeNotFoundError
	}

	return err.Error() == "not found"
}

// NilOrNotFound returns TRUE if the error is nil or a 404 / Not Found error.
// All other errors return FALSE
func NilOrNotFound(err error) bool {

	if isNil(err) {
		return true
	}

	if NotFound(err) {
		return true
	}

	return false
}

// IsInformational returns TRUE if the error `Code` is a 1xx / Informational error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#1xx_informational_response
func IsInformational(err error) bool {
	code := ErrorCode(err)
	return code >= 100 && code < 200
}

// IsSuccess returns TRUE if the error `Code` is a 2xx / Success error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#2xx_success
func IsSuccess(err error) bool {
	code := ErrorCode(err)
	return code >= 200 && code < 300
}

// IsRedirection returns TRUE if the error `Code` is a 3xx / Redirection error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#3xx_redirection
func IsRedirection(err error) bool {
	code := ErrorCode(err)
	return code >= 300 && code < 400
}

// IsClientError returns TRUE if the error `Code` is a 4xx / Client Error error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#4xx_client_errors
func IsClientError(err error) bool {
	code := ErrorCode(err)
	return code >= 400 && code < 500
}

// IsServerError returns TRUE if the error `Code` is a 5xx / Server Error error.
// https://en.wikipedia.org/wiki/List_of_HTTP_status_codes#5xx_server_errors
func IsServerError(err error) bool {
	code := ErrorCode(err)
	return code >= 500 && code < 600
}
