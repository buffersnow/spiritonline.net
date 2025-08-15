package web

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
)

type HttpError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Details struct {
	code    int      `json:"-"`
	Message string   `json:"message"`
	Err     error    `json:"-"`
	Context echo.Map `json:"context,omitempty"`
}

type baseTransmittableError struct {
	Error   string   `json:"error"`
	Message string   `json:"message"`
	Context echo.Map `json:"context,omitempty"`
}

func (e *Details) Error() string {
	return e.Err.Error()
}

func (e *Details) Unwrap() error {
	return e.Err
}

func (b Details) MarshalJSON() ([]byte, error) {
	return json.Marshal(baseTransmittableError{
		Error:   b.Err.Error(),
		Message: b.Message,
		Context: b.Context,
	})
}

// 400 - Bad HTTP Request
func BadRequestError(details *Details) error {
	details.code = http.StatusBadRequest
	return details
}

// 401 - Unauthorized Request
func AuthorizationError(details *Details) error {
	details.code = http.StatusUnauthorized
	return details
}

// 403 - Forbidden
func ForbiddenError(details *Details) error {
	details.code = http.StatusForbidden
	return details
}

// 404 - Not Found
func NotFoundError(details *Details) error {
	details.code = http.StatusNotFound
	return details
}

// 405 - Method Not Allowed
func BadMethodError(details *Details) error {
	details.code = http.StatusMethodNotAllowed
	return details
}

// 415 - Unsupported Media Type
func BadMediaTypeError(details *Details) error {
	details.code = http.StatusUnsupportedMediaType
	return details
}

// 451 - Unavailable For Legal Reason
func DMCAError(details *Details) error {
	details.code = http.StatusUnavailableForLegalReasons
	return details
}

// 500 - Internal Server Error
func InternalServerError(details *Details) error {
	details.code = http.StatusInternalServerError
	return details
}

func ErrorHandler(err error, c echo.Context) {
	var be *Details
	if errors.As(err, &be) {
		c.JSON(be.code, be)
		return
	}

	// Fallback for unexpected errors
	c.JSON(http.StatusTeapot, map[string]string{
		"message": "how did we get here?",
		"error":   err.Error(),
	})
}
