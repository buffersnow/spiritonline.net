package controllers

import (
	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"github.com/gofiber/fiber/v2"
)

func Download(c *fiber.Ctx) error {
	handlers := map[string]fiber.Handler{
		"count":    func(*fiber.Ctx) error { return nil },
		"list":     func(*fiber.Ctx) error { return nil },
		"contents": func(*fiber.Ctx) error { return nil },
	}

	action := c.FormValue("action")
	if len(action) == 0 {
		return protocol.NASReply(c, fiber.Map{
			"returncd": protocol.ReCD_InvalidAction,
		})
	}

	handler, ok := handlers[action]
	if !ok {
		return protocol.NASReply(c, fiber.Map{
			"returncd": protocol.ReCD_InvalidAction,
		})
	}

	return handler(c)
}
