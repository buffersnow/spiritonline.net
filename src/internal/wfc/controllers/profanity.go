package controllers

import "github.com/gofiber/fiber/v2"

//@ TODO: Implement (luxploit)
func Profanity(c *fiber.Ctx) error {
	return c.Status(fiber.StatusNotImplemented).SendString("arf!")
}
