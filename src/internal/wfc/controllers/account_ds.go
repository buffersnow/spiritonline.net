package controllers

import "github.com/gofiber/fiber/v2"

//@ TODO: Test with an actual DS (luxploit)

func AccountDS(c *fiber.Ctx) error {
	return c.SendString("arf!")
}
