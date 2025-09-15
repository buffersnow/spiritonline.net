package protocol

import (
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

//$ https://github.com/WiiLink24/wfc-server/blob/main/nas/auth.go#L410

func GetDateTime() string {
	return time.Now().Format("20060102150405")
}

func MarioKartOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		opt, err := red.Locate[settings.Options]()
		if err != nil {
			return web.BadLocateError(c, fmt.Errorf("wfc: protocol: %w", err))
		}

		if !opt.Service.Features["wfc_nas_mkwii_only"] {
			c.Next()
		}

		gamecd := c.Get("HTTP_X_GAMECD")
		if len(gamecd) == 0 {
			gamecd = c.FormValue("gamecd")
		}

		if len(gamecd) == 0 {
			return NASReply(c, fiber.Map{
				"returncd": ReCD_UnsupportedGame,
			})
		}

		if gamecd != "RMCP" {
			return NASReply(c, fiber.Map{
				"returncd": ReCD_UnsupportedGame,
			})
		}

		return c.Next()
	}
}
