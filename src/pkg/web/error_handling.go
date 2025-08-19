package web

import (
	"encoding/json"
	"errors"
	"fmt"

	"buffersnow.com/spiritonline/pkg/log"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

type HttpError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type Details struct {
	code    int       `json:"-"`
	Message string    `json:"message"`
	Err     error     `json:"-"`
	Context fiber.Map `json:"context,omitempty"`
}

type baseTransmittableError struct {
	Error   string    `json:"error"`
	Message string    `json:"message"`
	Context fiber.Map `json:"context,omitempty"`
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
func BadRequestError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusBadRequest
	c.Status(fiber.StatusBadRequest)
	return details
}

// 401 - Unauthorized Request
func AuthorizationError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusUnauthorized
	c.Status(fiber.StatusUnauthorized)
	return details
}

// 403 - Forbidden
func ForbiddenError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusForbidden
	c.Status(fiber.StatusForbidden)
	return details
}

// 404 - Not Found
func NotFoundError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusNotFound
	c.Status(fiber.StatusNotFound)
	return details
}

// 405 - Method Not Allowed
func BadMethodError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusMethodNotAllowed
	c.Status(fiber.StatusMethodNotAllowed)
	return details
}

// 415 - Unsupported Media Type
func BadMediaTypeError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusUnsupportedMediaType
	c.Status(fiber.StatusUnsupportedMediaType)
	return details
}

// 451 - Unavailable For Legal Reason
func DMCAError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusUnavailableForLegalReasons
	c.Status(fiber.StatusUnavailableForLegalReasons)
	return details
}

// 500 - Internal Server Error
func InternalServerError(c *fiber.Ctx, details *Details) error {
	details.code = fiber.StatusInternalServerError
	c.Status(fiber.StatusInternalServerError)
	return details
}

func ErrorHandler(c *fiber.Ctx, inerr error) error {
	logger, err := red.Locate[log.Logger]()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "bad service location",
			"error":   fmt.Errorf("web: %w", err).Error(),
		})
	}

	var be *Details
	if errors.As(inerr, &be) {
		js, err := json.Marshal(be.Context)
		if err != nil {
			logger.Error("HTTP Error Handler", "<IP: %s> Fallback! Error: %s", c.IP(), inerr.Error())
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"message": "bad json read",
				"error":   fmt.Errorf("web: %w", err).Error(),
			})
		}

		logger.Error("HTTP Error Handler",
			"<IP: %s> Message: %s, Context: %s, Error: %s",
			c.IP(), be.Message, string(js), be.Err.Error(),
		)

		return c.Status(be.code).JSON(be)
	}

	// Fallback for unexpected errors
	logger.Error("HTTP Error Handler", "<IP: %s> Fallback! Error: %s", c.IP(), inerr.Error())
	return c.Status(fiber.StatusTeapot).JSON(fiber.Map{
		"message": "how did we get here?",
		"error":   inerr.Error(),
	})
}
