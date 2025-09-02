package controllers

import (
	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"github.com/gofiber/fiber/v2"
)

//@ TODO: Test with an actual DS (luxploit)
//@ TODO: Actually implement, return 20107 atm (luxploit)

func AccountDS(c *fiber.Ctx) error {
	return protocol.NASReply(c, fiber.Map{
		"returncd": protocol.ReCD_UnsupportedGame,
	})
}
