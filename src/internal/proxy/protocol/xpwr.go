package protocol

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/version"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

func XPoweredBy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bld, err := red.Locate[version.BuildTag]()
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad service location",
				Err:     fmt.Errorf("web: %w", err),
			})
		}

		c.Set("WebProxy", fmt.Sprintf(
			"SpiritOnline/%s/%s (%s)", bld.GetVersion(), bld.GetService(), bld.GetConfig(),
		))

		return c.Next() // proceed to next middleware or handler
	}
}
