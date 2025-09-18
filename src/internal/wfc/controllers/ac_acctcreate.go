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
	//%   * macadr -> Console MAC Address (usually wireless nic)
	//%
	//% fiber.Ctx:
	//%   * c.IP() -> Console IP Address

	repo, err := red.Locate[repositories.WFCRepo]()
	if err != nil {
		return web.BadLocateError(c, fmt.Errorf("wfc: controllers: %w", err))
	}

	wfcid, err := protocol.GetWFCAccountID(repo, repositories.WFCAccountQuery{
		Serial: c.FormValue("csnum"),
		FC:     cast.ToInt64(c.FormValue("cfc")),
		IP:     c.IP(),
		MAC:    c.FormValue("macadr"),
	})
	if err != nil {
		return web.InternalServerError(c, &web.Details{
			Message: "bad db query",
			Err:     fmt.Errorf("wfc: controllers: %w", err),
		})
	}

	suspension, err := repo.Suspension.Get(wfcid)
	if err != nil && err != sql.ErrNoRows {
		return web.InternalServerError(c, &web.Details{
			Message: "bad db query",
			Err:     fmt.Errorf("wfc: controllers: %w", err),
		})
	} else if err == sql.ErrNoRows {
		return protocol.NASReply(c, fiber.Map{
			"returncd": protocol.ReCD_AccountCreate,
			"userid":   wfcid,
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
		"userid":   wfcid,
	})

}
