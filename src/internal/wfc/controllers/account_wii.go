package controllers

import "github.com/gofiber/fiber/v2"

func AccountWii(c *fiber.Ctx) error {
	return c.SendString("arf!")
}
