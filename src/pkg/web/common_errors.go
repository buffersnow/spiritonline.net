package web

import "github.com/gofiber/fiber/v2"

// Prepared [InternalServerError] for bad service locations
func BadLocateError(c *fiber.Ctx, err error) error {
	return InternalServerError(c, &Details{
		Message: "bad service location",
		Err:     err,
	})
}
