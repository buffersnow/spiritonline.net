package controllers

import (
	"database/sql"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"buffersnow.com/spiritonline/internal/wfc/repositories"
	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
	"github.com/spf13/cast"
)

func AC_Login(c *fiber.Ctx) error {

	//% refer to ac_acctcreate.go for which factors determine a user

	repo, err := red.Locate[repositories.WFCRepo]()
	if err != nil {
		return web.BadLocateError(c, fmt.Errorf("wfc: controllers: %w", err))
	}

	query := repositories.WFCAccountQuery{
		Serial: c.FormValue("csnum"),
		FC:     cast.ToInt64(c.FormValue("cfc")),
		IP:     c.IP(),
		MAC:    c.FormValue("macadr"),
	}

	wfcid, err := protocol.GetWFCAccountID(repo, query)
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
	} else if suspension.AuditID != 0 /*should be valid*/ {
		if !suspension.BanExpiresOn.Valid {
			return protocol.NASReply(c, fiber.Map{
				"returncd": protocol.ReCD_BannedFromWFC,
			})
		} else if suspension.BanExpiresOn.Time.Before(time.Now()) {
			return protocol.NASReply(c, fiber.Map{
				"returncd": protocol.ReCD_TempBannedFromWFC,
			})
		}

	}

	challenge := util.RandomString(8)
	token, err := gp.PickleWFCToken(gp.WFCAuthToken{
		WFCID:     wfcid,
		GameCode:  c.FormValue("gamecd"),
		RegionID:  util.HexToByte(c.FormValue("region")),
		Serial:    query.Serial,
		ConsoleFC: query.FC,
		MAC:       query.MAC,
		IP:        query.IP,
		Challenge: challenge,
		UnitCode:  cast.ToInt8(c.FormValue("unitcd")),
	})

	if err != nil {
		return web.InternalServerError(c, &web.Details{
			Message: "bad token generation",
			Err:     fmt.Errorf("wfc: controllers: %w", err),
		})
	}

	return protocol.NASReply(c, fiber.Map{
		"returncd":  protocol.ReCD_Login,
		"challenge": challenge,
		"token":     token,
	})
}
