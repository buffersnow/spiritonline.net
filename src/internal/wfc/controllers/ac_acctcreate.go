package controllers

import (
	"database/sql"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"buffersnow.com/spiritonline/internal/wfc/repositories"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
	"github.com/spf13/cast"
)

//$ https://github.com/insanekartwii/wfc-server/blob/main/nas/auth.go#L179

func AC_AccountCreate(c *fiber.Ctx) error {

	//% Fields we use to determine a User
	//% Post Body:
	//%   * csnum  -> Console Serial Number (aka CiD)
	//%   * cfc    -> Console NandID
	//%   * macadr -> Console MAC Address (usually wireless nic)
	//%
	//% fiber.Ctx:
	//%   * c.IP() -> Console IP Address

	repo, err := red.Locate[repositories.WFCRepo]()
	if err != nil {
		return web.BadLocateError(c, fmt.Errorf("wfc: controllers: %w", err))
	}

	cfc, err := cast.ToInt64E(c.FormValue("cfc"))
	if err != nil {
		return web.BadRequestError(c, &web.Details{
			Message: "invalid nandid/cfc",
			Err:     fmt.Errorf("wfc: controllers: cast: %w", err),
		})
	}

	query := repositories.WFCAccountQuery{
		ConsoleID: c.FormValue("csnum"),
		NandID:    cfc,
		IP:        c.IP(),
		MAC:       c.FormValue("macadr"),
	}

	acc, err := repo.Account.Get(query)
	if err != nil && err != sql.ErrNoRows {
		return web.BadRequestError(c, &web.Details{
			Message: "bad db query",
			Err:     fmt.Errorf("wfc: controllers: %w", err),
		})
	} else if err == sql.ErrNoRows {
		wfcid, err := repo.Account.Insert(query)
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad db insert",
				Err:     fmt.Errorf("wfc: controllers: %w", err),
			})
		}

		return protocol.NASReply(c, fiber.Map{
			"returncd": protocol.ReCD_AccountCreate,
			"userid":   wfcid,
		})
	}

	suspension, err := repo.Suspension.Get(acc.WFCID)
	if err != nil && err != sql.ErrNoRows {
		return web.InternalServerError(c, &web.Details{
			Message: "bad db query",
			Err:     fmt.Errorf("wfc: controllers: %w", err),
		})
	} else if err == sql.ErrNoRows {
		return protocol.NASReply(c, fiber.Map{
			"returncd": protocol.ReCD_AccountCreate,
			"userid":   acc.WFCID,
		})
	}

	if !suspension.BanExpiresOn.Valid {
		return protocol.NASReply(c, fiber.Map{
			"returncd": protocol.ReCD_BannedFromWFC,
		})
	} else if suspension.BanExpiresOn.Time.Before(time.Now()) {
		return protocol.NASReply(c, fiber.Map{
			"returncd": protocol.ReCD_TempBannedFromWFC,
		})
	}

	return protocol.NASReply(c, fiber.Map{
		"returncd": protocol.ReCD_AccountCreate,
		"userid":   acc.WFCID,
	})

}
