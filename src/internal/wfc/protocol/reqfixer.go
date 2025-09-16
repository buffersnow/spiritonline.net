package protocol

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

func RequestFixup() fiber.Handler {
	return func(c *fiber.Ctx) error {

		//& validated in protocol/validation.go
		i64_unitcd := cast.ToInt64(c.FormValue("unitcd"))

		str_lang := c.FormValue("lang", "FE")
		if u8_lang := util.HexToByte(str_lang); u8_lang == 0xFF || u8_lang == 0x00 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid lang",
				Err:     fmt.Errorf("wfc: protocol: lang was outside of boundaries"),
			})
		}

		if i64_unitcd == UnitCD_NintendoWii {
			//$ https://github.com/Retro-Rewind-Team/wfc-server/blob/main/nas/auth.go#L318
			str_region := c.FormValue("region", "FE")
			if u8_region := util.HexToByte(str_region); u8_region == 0xFF || u8_region == 0x00 {
				return web.BadRequestError(c, &web.Details{
					Message: "invalid region",
					Err:     fmt.Errorf("wfc: protocol: region was outside of boundaries"),
				})
			}
		}

		str_gamecd := c.FormValue("gamecd")
		if len(str_gamecd) != 4 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid gamecd",
				Err:     fmt.Errorf("wfc: protocol: length was less then or greater then 4"),
			})
		}

		//$ https://github.com/Retro-Rewind-Team/wfc-server/blob/main/nas/auth.go#L223
		str_gsbrcd := c.FormValue("gsbrcd")
		if len(str_gsbrcd) < 4 {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid gsbrcd",
				Err:     fmt.Errorf("wfc: protocol: length of gsbrcd is 0"),
			})
		}

		//& WiiLink just sets gsbrcd to first 3 chars of gamecd + "J" but
		//& doesn't include the random bits (which are in my test always 8)
		//? so hopefully this doesn't blow anything up, guess we'll find out
		if len(str_gsbrcd) == 0 {
			str_gsbrcd = str_gamecd[:3] + "J"
			str_gsbrcd += util.RandomString(8)
		} else if str_gsbrcd[:3] != str_gamecd[:3] {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid gsbrcd",
				Err:     fmt.Errorf("wfc: protocol: mismatch between gsbrcd and gamecd"),
			})
		}

		return c.Next()
	}
}
