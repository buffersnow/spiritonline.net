package protocol

import (
	"fmt"
	"net/url"
	"strings"

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
		return web.InternalServerError(c, &web.Details{
			Message: "bad service location",
			Err:     fmt.Errorf("wfc: protocol: %w", err),
		})
	}

	urlVals := url.Values{}
	{ // these values are always present in responses (or should be!)
		urlVals.Set("retry", "0")
		urlVals.Set("datetime", GetDateTime())
		urlVals.Set("locator", "gamespy.com")
	}

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
			num, err := cast.ToIntE(str)
			if err != nil {
				return web.InternalServerError(c, &web.Details{
					Message: "bad returncd value",
					Err:     fmt.Errorf("wfc: protocol: %w", err),
				})
			}
			str = fmt.Sprintf("%03d", num)
		}

		b64, err := sec.Encoding.EncodeB64_Wii([]byte(str))
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad b64_wii encoding",
				Err:     fmt.Errorf("wfc: protocol: %w", err),
			})
		}

		urlVals.Set(k, string(b64))
	}

	resp := urlVals.Encode()
	resp = strings.ReplaceAll(resp, "%2A", "*")

	return c.Type(fiber.MIMEApplicationForm).SendString(resp)
}
