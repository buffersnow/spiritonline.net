package controllers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"buffersnow.com/spiritonline/internal/wfc/repositories"
	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
	"github.com/spf13/cast"
)

func AC_ServiceLocate(c *fiber.Ctx) error {

	//$ https://github.com/barronwaffles/dwc_network_server_emulator/blob/master/nas_server.py#L153
	svc := c.FormValue("svc")
	if len(svc) != 0 {
		_, err := cast.ToIntE(svc)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid svc",
				Err:     fmt.Errorf("wfc: protocol: cast: %w", err),
			})
		}
	}

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

	switch svc {
	case "9000":
		return protocol.NASReply(c, fiber.Map{
			"returncd":   protocol.ReCD_ServiceLocate,
			"statusdata": "Y",
			"svchost":    "dls1.nintendowifi.net",
			"token":      token,
		})
	case "9001":
		return protocol.NASReply(c, fiber.Map{
			"returncd":     protocol.ReCD_ServiceLocate,
			"statusdata":   "Y",
			"svchost":      "dls1.nintendowifi.net",
			"servicetoken": token,
		})
	default: //& covers "0000" requested by Pokemon GTS and empty by Boom Street
		return protocol.NASReply(c, fiber.Map{
			"returncd":     protocol.ReCD_ServiceLocate,
			"statusdata":   "Y",
			"svchost":      "n/a",
			"servicetoken": token,
		})
	}
}
