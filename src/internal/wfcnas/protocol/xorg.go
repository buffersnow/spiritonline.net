package protocol

import "github.com/labstack/echo/v4"

func XOrganization(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		c.Response().Header().Set("X-Organization", "Nintendo")
		return next(c)
	}
}
