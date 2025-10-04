package protocol

import (
	"fmt"
	"net/url"
	"strings"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/security"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
	"github.com/spf13/cast"
)

//$ https://github.com/WiiLink24/wfc-server/blob/main/nas/auth.go#L162

func NASReply(c *fiber.Ctx, params fiber.Map) error {

	sec, err := red.Locate[security.Security]()
	if err != nil {
		return web.BadLocateError(c, fmt.Errorf("wfc: protocol: %w", err))
	}

	logger, err := red.Locate[log.Logger]()
	if err != nil {
		return web.BadLocateError(c, fmt.Errorf("wfc: protcol: %w", err))
	}

	{ // these values are always present in responses (or should be!)
		params["retry"] = 0
		params["datetime"] = GetDateTime()
		params["locator"] = "gamespy.com"
	}

	urlVals := url.Values{}
	logVals := url.Values{}
	recdnum := 0
	for k, v := range params {

		str, err := cast.ToStringE(v)
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad response string",
				Err:     fmt.Errorf("wfc: protocol: %w", err),
			})
		}

		//& Special case: pad returncd to 3 digits
		if k == "returncd" {
			recdnum, err = cast.ToIntE(str)
			if err != nil {
				return web.InternalServerError(c, &web.Details{
					Message: "bad returncd value",
					Err:     fmt.Errorf("wfc: protocol: %w", err),
				})
			}
			str = fmt.Sprintf("%03d", recdnum)
		}

		b64, err := sec.Encoding.EncodeB64_Wii([]byte(str))
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad b64_wii encoding",
				Err:     fmt.Errorf("wfc: protocol: %w", err),
			})
		}

		urlVals.Set(k, string(b64))
		logVals.Set(k, str)
	}

	resp := urlVals.Encode()
	resp = strings.ReplaceAll(resp, "%2A", "*")

	logResp := logVals.Encode()
	logResp = strings.ReplaceAll(logResp, "%2A", "*")

	unitcd := cast.ToInt(c.FormValue("unitcd"))

	logger.Trace(log.DEBUG_SERVICE, "WFC NAS Response",
		"<IP: %s> %s %s! (ReCD: %03d) Sent with data: %s",
		c.IP(), GetEndpoint(c), GetReCDMeaning(recdnum, unitcd), recdnum, logResp,
	)

	return c.Type(fiber.MIMEApplicationForm).Status(200).SendString(resp)
}
