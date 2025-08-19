package controllers

import "github.com/gofiber/fiber/v2"

func NasTest(c *fiber.Ctx) error {
	return c.Render("nastest", nil)
}
