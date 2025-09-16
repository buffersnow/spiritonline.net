package protocol

import (
	"encoding/hex"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

// Validates unitcd, userid, sdkver, makercd, macadr, devtime, cfc and csnum
// coming from any POST request issued by a Wii or DS.
func ValidateRequest() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get(fiber.HeaderContentType) != fiber.MIMEApplicationForm {
			return c.Next()
		}

		if c.Method() != fiber.MethodPost {
			return c.Next()
		}

		//@ TODO: Make sure this work on the DS! Currently only works tested on Wii

		//% action validation is done in controllers/account_[wii/ds].go
		//% gsbrcd, lang, gamecd and region are validated in protocol/reqfixer.go
		//% ingamesn and devname profanity is validated in protocol/ac_[login/acctcreate].go

		//& its important to this one first
		//& used for per-console behaviour
		str_unitcd := c.FormValue("unitcd")
		if len(str_unitcd) == 0 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid unitcd",
				Err:     fmt.Errorf("wfc: protocol: length of unitcd is 0"),
			})
		}

		i64_unitcd, err := cast.ToInt64E(str_unitcd)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid unitcd",
				Err:     fmt.Errorf("wfc: protocol: cast: %w", err),
			})
		}

		str_userid := c.FormValue("userid")
		if len(str_userid) == 0 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid userid",
				Err:     fmt.Errorf("wfc: protocol: length of userid is 0"),
			})
		}

		i64_userid, err := cast.ToInt64E(str_userid)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid userid",
				Err:     fmt.Errorf("wfc: protocol: cast: %w", err),
			})
		}

		//$ https://github.com/Retro-Rewind-Team/wfc-server/blob/main/nas/auth.go#L210
		if i64_userid >= 0x80000000000 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid userid size",
				Err:     fmt.Errorf("wfc: protocol: size of userid exceeded maximum"),
			})
		}

		str_sdkver := c.FormValue("sdkver")
		if len(str_sdkver) == 0 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid sdkver",
				Err:     fmt.Errorf("wfc: protocol: length of sdkver is 0"),
			})
		}

		i64_sdkver, err := cast.ToInt64E(str_sdkver)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid sdkver",
				Err:     fmt.Errorf("wfc: protocol: cast: %w", err),
			})
		}

		//$ http://www.pipian.com/ierukana/hacking/ds_nwc.html#example_nas_new_account
		//& sdkver is formated as %03d%03d in the SDK and as we
		//& assume that there's no SDK that uses major version
		//& less then 1 (formated as 001). but because of a bug
		//& in cast, 001000 is read in as 512 (probably octal)
		//& so 512 is the lowest value, not 1000.. thanks spf13
		if i64_sdkver < 512 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid sdkver",
				Err:     fmt.Errorf("wfc: protocol: major version of sdkver was less then 001"),
			})
		}

		str_makercd := c.FormValue("makercd")
		if len(str_makercd) == 0 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid sdkver",
				Err:     fmt.Errorf("wfc: protocol: length of makercd is 0"),
			})
		}

		i64_makercd, err := cast.ToInt64E(str_makercd)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid makercd",
				Err:     fmt.Errorf("wfc: protocol: cast: %w", err),
			})
		}

		//$ http://www.pipian.com/ierukana/hacking/ds_nwc.html#example_nas_new_account
		//& makercd is formated as %02d, so there shouldn't be a value bigger then 99
		//& and there also cannot be a value of 0 as nintendo uses 01, so thats invalid
		if i64_makercd <= 0 || i64_makercd > 99 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid makercd",
				Err:     fmt.Errorf("wfc: protocol: makercd value was greater then max and lt 0"),
			})
		}

		//$ https://standards.ieee.org/wp-content/uploads/import/documents/tutorials/eui.pdf
		//& MAC addresses can only be 48 bits ever
		//& hence 6 bytes which makes 12 digits
		str_macaddress := c.FormValue("macadr")
		if len(str_macaddress) != 12 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid mac address",
				Err:     fmt.Errorf("wfc: protocol: length of mac address is invalid"),
			})
		}

		_, err = hex.DecodeString(str_macaddress)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid mac address",
				Err:     fmt.Errorf("wfc: protocol: %w", err),
			})
		}

		str_devtime := c.FormValue("devtime")
		if len(str_devtime) == 0 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid sdkver",
				Err:     fmt.Errorf("wfc: protocol: length of devtime is 0"),
			})
		}

		devtime, err := time.Parse("060102150405", str_devtime)
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid devtime",
				Err:     fmt.Errorf("wfc: protocol: time: %w", err),
			})
		}

		if !util.WithinTimezoneDrift(devtime) {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid devtime",
				Err:     fmt.Errorf("wfc: protocol: devtime exceeded maximum timezone drift"),
			})
		}

		if i64_unitcd == UnitCD_NintendoWii {
			//& CFC = NandID(?) so if this null its perfectly fine (i think?)
			str_cfc := c.FormValue("cfc")
			_, err = cast.ToInt64E(str_cfc)
			if err != nil {
				return web.BadRequestError(c, &web.Details{
					Message: "invalid makercd",
					Err:     fmt.Errorf("wfc: protocol: cast: %w", err),
				})
			}
		}

		str_csnum := c.FormValue("csnum")
		if len(str_csnum) == 0 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid csnum",
				Err:     fmt.Errorf("wfc: protocol: length of csnum is 0"),
			})
		}

		//$ https://github.com/Retro-Rewind-Team/wfc-server/blob/main/nas/auth.go#L275
		if len(str_csnum) <= 8 || len(str_csnum) > 16 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid csnum",
				Err:     fmt.Errorf("wfc: protocol: length of csnum is 0"),
			})
		}

		return c.Next()
	}
}
