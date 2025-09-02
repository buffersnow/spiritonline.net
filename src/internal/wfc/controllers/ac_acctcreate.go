package controllers

import (
	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"github.com/gofiber/fiber/v2"
)

//$ https://github.com/insanekartwii/wfc-server/blob/main/nas/auth.go#L179

func AC_AccountCreate(c *fiber.Ctx) error {
	return protocol.NASReply(c, fiber.Map{
		"returncd": protocol.ReCD_AccountCreate,
	})
}
