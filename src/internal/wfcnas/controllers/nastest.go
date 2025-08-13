package controllers

import "github.com/labstack/echo/v4"

func NasTest(c echo.Context) error {
	return c.Render(200, "nastest", nil)
}
