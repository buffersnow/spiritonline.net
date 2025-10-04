package controllers

import "github.com/gofiber/fiber/v2"

func Test(c *fiber.Ctx) error {
	return c.Render("wfc/nastest", nil)
}
