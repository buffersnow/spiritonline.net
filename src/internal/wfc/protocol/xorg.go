package protocol

import (
	"github.com/gofiber/fiber/v2"
)

func XOrganization() fiber.Handler {
	return func(c *fiber.Ctx) error {
		c.Set("X-Organization", "Nintendo")
		return c.Next()
	}
}
