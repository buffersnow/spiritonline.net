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

	//$ https://github.com/barronwaffles/dwc_network_server_emulator/blob/master/dls1_server.py#L79
	//$ https://github.com/Retro-Rewind-Team/wfc-server/blob/main/nas/auth.go#L159
	//? Both wwfc and AltWFC both put this header here, so it's probably a good idea to set it
	c.Set("X-DLS-Host", "http://127.0.0.1")

	return handler(c)
}
