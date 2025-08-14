package controllers

import "github.com/labstack/echo/v4"

func Account(c echo.Context) error {
	return c.String(200, "arf!")
}
