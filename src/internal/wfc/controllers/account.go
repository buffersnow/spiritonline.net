package controllers

import "github.com/gofiber/fiber/v2"

func Account(c *fiber.Ctx) error {
	return c.SendString("arf!")
}
