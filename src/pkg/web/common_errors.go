package web

import "github.com/gofiber/fiber/v2"

// Prepared [InternalServerError] for bad service locations
func BadLocateError(err error) error {
	return InternalServerError(&Details{
		Message: "bad service location attempt",
		Err:     err,
	})
}

// 400 - Bad HTTP Request
func BadRequestError(details *Details) error {
	details.Code = fiber.StatusBadRequest
	return details
}

// 401 - Unauthorized Request
func AuthorizationError(details *Details) error {
	details.Code = fiber.StatusUnauthorized
	return details
}

// 403 - Forbidden
func ForbiddenError(details *Details) error {
	details.Code = fiber.StatusForbidden
	return details
}

// 404 - Not Found
func NotFoundError(details *Details) error {
	details.Code = fiber.StatusNotFound
	return details
}

// 405 - Method Not Allowed
func BadMethodError(details *Details) error {
	details.Code = fiber.StatusMethodNotAllowed
	return details
}

// 415 - Unsupported Media Type
func BadMediaTypeError(details *Details) error {
	details.Code = fiber.StatusUnsupportedMediaType
	return details
}

// 451 - Unavailable For Legal Reason
func DMCAError(details *Details) error {
	details.Code = fiber.StatusUnavailableForLegalReasons
	return details
}

// 500 - Internal Server Error
func InternalServerError(details *Details) error {
	details.Code = fiber.StatusInternalServerError
	return details
}

// 502 - Bad Gateway Error
func BadGatewayError(details *Details) error {
	details.Code = fiber.StatusBadGateway
	return details
}
