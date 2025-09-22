package protocol

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
	"golang.org/x/text/encoding/unicode"
)

func RequestFixup() fiber.Handler {
	return func(c *fiber.Ctx) error {

		//& validated in protocol/validation.go
		i64_unitcd := cast.ToInt64(c.FormValue("unitcd"))

		str_lang := c.FormValue("lang", "FE")
		if u8_lang := util.HexToByte(str_lang); u8_lang == 0xFF {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid lang",
				Err:     fmt.Errorf("wfc: protocol: lang was outside of boundaries"),
			})
		}

		if i64_unitcd == UnitCD_NintendoWii {
			//$ https://github.com/Retro-Rewind-Team/wfc-server/blob/main/nas/auth.go#L318
			str_region := c.FormValue("region", "FE")
			if u8_region := util.HexToByte(str_region); u8_region == 0xFF {
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
				Err:     fmt.Errorf("wfc: protocol: length of gsbrcd less then the required minimum"),
			})
		}

		//& WiiLink just sets gsbrcd to first 3 chars of gamecd + "J" but
		//& doesn't include the random bits (which are in my test always 8)
		//? so hopefully this doesn't blow anything up, guess we'll find out
		if len(str_gsbrcd) == 0 {
			str_gsbrcd = str_gamecd[:3] + "J"
			str_gsbrcd += util.RandomString(8)
			c.Request().PostArgs().Set("gsbrcd", str_gsbrcd)
		} else if str_gsbrcd[:3] != str_gamecd[:3] {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid gsbrcd",
				Err:     fmt.Errorf("wfc: protocol: mismatch between gsbrcd and gamecd"),
			})
		}

		//$ https://github.com/insanekartwii/wfc-server/blob/main/nas/auth.go#L64
		ingamesn := c.FormValue("ingamesn")
		devname := c.FormValue("devname")
		words := c.FormValue("words")

		endian := unicode.BigEndian

		if i64_unitcd == UnitCD_NintendoDS {
			endian = unicode.LittleEndian

			if len(devname) == 0 {
				return web.BadRequestError(c, &web.Details{
					Message: "invalid devname",
					Err:     fmt.Errorf("wfc: protocol: length of devname is 0"),
				})
			}

			utf, err := util.FromUTF16(unicode.LittleEndian, []byte(devname))
			if err != nil {
				return web.BadRequestError(c, &web.Details{
					Message: "invalid devname",
					Err:     fmt.Errorf("wfc: protocol: %w", err),
				})
			}

			c.Request().PostArgs().Set("devname", utf)
		}

		if len(ingamesn) == 0 && i64_unitcd == UnitCD_NintendoWii {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid ingamesn",
				Err:     fmt.Errorf("wfc: protocol: length of ingamesn is 0"),
			})
		}

		if len(ingamesn) != 0 { //& ingamesn is not required on DS
			utf, err := util.FromUTF16(endian, []byte(ingamesn))
			if err != nil {
				return web.BadRequestError(c, &web.Details{
					Message: "invalid ingamesn",
					Err:     fmt.Errorf("wfc: protocol: %w", err),
				})
			}

			c.Request().PostArgs().Set("ingamesn", utf)
		}

		if len(words) != 0 {
			utf, err := util.FromUTF16(endian, []byte(words))
			if err != nil {
				return web.BadRequestError(c, &web.Details{
					Message: "invalid profanity words",
					Err:     fmt.Errorf("wfc: protocol: %w", err),
				})
			}

			c.Request().PostArgs().Set("words", utf)
		}

		return c.Next()
	}
}
