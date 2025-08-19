package controllers

import "github.com/gofiber/fiber/v2"

func Product(c *fiber.Ctx) error {
	return c.SendString("arf!")
}
