package protocol

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/settings"
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

		opt, err := red.Locate[settings.Options]()
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad service location",
				Err:     fmt.Errorf("web: %w", err),
			})
		}

		c.Request().Header.Set("X-Web-Proxy", fmt.Sprintf(
			"SpiritOnline/%s/%s (%s)", bld.GetVersion(), bld.GetService(), bld.GetConfig(),
		))
		c.Request().Header.Set("X-Proxy-Tag", opt.Spirit.ServiceTag)

		return c.Next() // proceed to next middleware or handler
	}
}
