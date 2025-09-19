package controllers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

//@ TODO: Test with an actual DS (luxploit)
//@ TODO: Actually implement, return 20107 atm (luxploit)

func AccountDS(c *fiber.Ctx) error {
	opt, err := red.Locate[settings.Options]()
	if err != nil {
		return web.BadLocateError(c, fmt.Errorf("proxy: protocol: %w", err))
	}

	if opt.Service.Features["wfc_nas_enable_ds"] {
		panic("absolute cinema, don't touch stuff >:(")
	}

	return protocol.NASReply(c, fiber.Map{
		"returncd": protocol.ReCD_UnsupportedGame,
	})
}
