package protocol

import (
	"time"

	"github.com/gofiber/fiber/v2"
)

//$ https://github.com/WiiLink24/wfc-server/blob/main/nas/auth.go#L410

func GetDateTime() string {
	return time.Now().Format("20060102150405")
}

func MarioKartOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
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
