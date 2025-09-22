package protocol

import (
	"database/sql"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/internal/wfc/repositories"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

//$ https://github.com/WiiLink24/wfc-server/blob/main/nas/auth.go#L410

func GetDateTime() string {
	return time.Now().Format("060102150405")
}

func MarioKartOnly() fiber.Handler {
	return func(c *fiber.Ctx) error {
		opt, err := red.Locate[settings.Options]()
		if err != nil {
			return web.BadLocateError(c, fmt.Errorf("wfc: protocol: %w", err))
		}

		if !opt.Service.Features["wfc_nas_mkwii_only"] {
			return c.Next()
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

		if gamecd != "RMCP" && gamecd != "RMCE" && gamecd != "RMCJ" {
			return NASReply(c, fiber.Map{
				"returncd": ReCD_UnsupportedGame,
			})
		}

		return c.Next()
	}
}

func GetEndpoint(c *fiber.Ctx) string {
	if str := c.FormValue("action"); len(str) != 0 {
		return fmt.Sprintf("<Action: %s> <Endpoint: %s>", str, c.Path())
	}

	return fmt.Sprintf("<Endpoint: %s>", c.Path())
}

func GetWFCAccountID(repo *repositories.WFCRepo, query repositories.WFCAccountQuery) (int64, error) {
	wfcid, err := repo.Account.GetWFCID(query)

	if err != nil && err != sql.ErrNoRows {
		return 0, err
	}

	if err == sql.ErrNoRows {
		wfcid, err = repo.Account.Insert(query)
		if err != nil {
			return 0, err
		}

		return wfcid, nil
	}

	err = repo.Account.Update(query)
	if err != nil {
		return 0, err
	}

	return wfcid, nil
}
