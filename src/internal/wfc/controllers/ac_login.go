package controllers

import (
	"database/sql"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"buffersnow.com/spiritonline/internal/wfc/repositories"
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
	}

	if err == sql.ErrNoRows {
		wfcid, err := repo.Account.Insert(query)
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad db insert",
				Err:     fmt.Errorf("wfc: controllers: %w", err),
			})
		}

		acc, err = repo.Account.GetByWFCID(wfcid)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "bad db refresh",
				Err:     fmt.Errorf("wfc: controllers: %w", err),
			})
		}
	}

	suspension, err := repo.Suspension.Get(acc.WFCID)
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
	token, err := protocol.CreateToken(protocol.AuthToken{
		WFCID:     acc.WFCID,
		GameCode:  c.FormValue("gamecd"),
		RegionID:  util.HexToByte(c.FormValue("region")),
		ConsoleID: query.ConsoleID,
		NandID:    query.NandID,
		MAC:       query.MAC,
		IP:        query.IP,
		Challenge: challenge,
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
