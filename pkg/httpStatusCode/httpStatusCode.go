package httpStatusCode

type HTTP_STATUS_CODE int

const (
	OK_CODE                 HTTP_STATUS_CODE = 200
	NOT_MODIFIED_CODE       HTTP_STATUS_CODE = 304
	FORBIDDEN_CODE          HTTP_STATUS_CODE = 403
	NOT_FOUND_CODE          HTTP_STATUS_CODE = 404
	METHOD_NOT_ALLOWED_CODE HTTP_STATUS_CODE = 405
)

type HTTP_STATUS_TEXT string

const (
	OK_TEXT                 HTTP_STATUS_TEXT = "OK"
	NOT_MODIFIED_TEXT       HTTP_STATUS_TEXT = "Not Modified"
	FORBIDDEN_TEXT          HTTP_STATUS_TEXT = "Forbidden"
	NOT_FOUND_TEXT          HTTP_STATUS_TEXT = "Not Found"
	METHOD_NOT_ALLOWED_TEXT HTTP_STATUS_TEXT = "Method Not Allowed"
)
