package web

import (
	"time"

	"buffersnow.com/spiritonline/pkg/log"
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
