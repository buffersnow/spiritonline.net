package web

import (
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/version"
	"github.com/labstack/echo/v4"
)

func RequestLogging(prefix string, logger *log.Logger) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			start := time.Now()
			err := next(c)
			latency := time.Since(start)

			req := c.Request()
			res := c.Response()

			ip := c.RealIP()

			logger.Trace(log.DEBUG_SERVICE, prefix, "<IP: %s> <Time: %v> <Status: %d> %s %s", ip, latency, res.Status, req.Method, req.URL.RequestURI())

			return err
		}
	}
}

func XPoweredBy(bld *version.BuildTag) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			c.Response().Header().Set("x-powered-by", "buffersnow.com")
			c.Response().Header().Set("Server", fmt.Sprintf(
				"SpiritOnline/%s/%s (%s)", bld.GetVersion(), bld.GetService(), bld.GetConfig(),
			))
			return next(c)
		}
	}
}
