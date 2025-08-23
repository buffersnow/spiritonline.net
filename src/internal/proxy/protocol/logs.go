package protocol

import (
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

func RequestLogging() fiber.Handler {
	return func(c *fiber.Ctx) error {
		start := time.Now()

		logger, err := red.Locate[log.Logger]()
		if err != nil {
			return web.InternalServerError(c, &web.Details{
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
		target := c.Get("Host")
		host := c.Get("X-Forwarded-Host")

		if host == "" {
			host = "none"
		}

		logger.Trace(log.DEBUG_TRAFFIC, "HTTP Request Logger",
			"<IP: %s> <Time: %v> <Status: %d> <Target: %s> <Host: %s> %s %s",
			realip, latency, status, target, host, method, uri,
		)

		return err
	}
}
