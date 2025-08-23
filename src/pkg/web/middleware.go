package web

import (
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/version"

	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

func RequestLogging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		logger, err := red.Locate[log.Logger]()
		if err != nil {
			return InternalServerError(c, &Details{
				Message: "bad service location",
				Err:     fmt.Errorf("web: %w", err),
			})
		}

		err = c.Next()

		latency := time.Since(start)
		realip := c.IP()
		status := c.Response().StatusCode()
		method := c.Method()
		uri := c.OriginalURL()

		logger.Trace(log.DEBUG_TRAFFIC, "HTTP Request Logger",
			"<IP: %s> <Time: %v> <Status: %d> %s %s",
			realip, latency, status, method, uri,
		)

		return err
	}
}

func XPoweredBy() fiber.Handler {
	return func(c *fiber.Ctx) error {
		bld, err := red.Locate[version.BuildTag]()
		if err != nil {
			return InternalServerError(c, &Details{
				Message: "bad service location",
				Err:     fmt.Errorf("web: %w", err),
			})
		}

		opt, err := red.Locate[settings.Options]()
		if err != nil {
			return InternalServerError(c, &Details{
				Message: "bad service location",
				Err:     fmt.Errorf("web: %w", err),
			})
		}

		c.Set("X-Powered-By", "buffersnow.com")

		c.Set("X-Web-Proxy", c.Get("X-Web-Proxy"))
		c.Set("X-Proxy-Tag", c.Get("X-Proxy-Tag"))

		c.Set("X-Service-Tag", opt.Spirit.ServiceTag)
		c.Set("Server", fmt.Sprintf(
			"SpiritOnline/%s/%s (%s)", bld.GetVersion(), bld.GetService(), bld.GetConfig(),
		))

		return c.Next() // proceed to next middleware or handler
	}
}
