package web

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/version"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/luxploit/red"
)

func RequestLogging(prefix string) echo.MiddlewareFunc {
	return middleware.RequestLoggerWithConfig(middleware.RequestLoggerConfig{
		LogStatus:   true,
		LogURI:      true,
		LogMethod:   true,
		LogLatency:  true,
		LogRemoteIP: true,
		LogValuesFunc: func(c echo.Context, v middleware.RequestLoggerValues) error {
			logger, err := red.Locate[log.Logger]()
			if err != nil {
				return InternalServerError(&Details{
					Message: "bad service location",
					Err:     fmt.Errorf("web: %w", err),
				})
			}

			logger.Trace(log.DEBUG_SERVICE, prefix,
				"<IP: %s> <Time: %v> <Status: %d> %s %s",
				v.RemoteIP, v.Latency, v.Status, v.Method, v.URI,
			)

			return nil
		},
	})
}

func XPoweredBy(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		bld, err := red.Locate[version.BuildTag]()
		if err != nil {
			return InternalServerError(&Details{
				Message: "bad service location",
				Err:     fmt.Errorf("web: %w", err),
			})
		}

		c.Response().Header().Set("x-powered-by", "buffersnow.com")
		c.Response().Header().Set("Server", fmt.Sprintf(
			"SpiritOnline/%s/%s (%s)", bld.GetVersion(), bld.GetService(), bld.GetConfig(),
		))
		return next(c)
	}
}
