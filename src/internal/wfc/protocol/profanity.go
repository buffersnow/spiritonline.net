package protocol

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/web"
	goaway "github.com/TwiN/go-away"
	"github.com/gofiber/fiber/v2"
	"github.com/spf13/cast"
)

func ProfanityFilter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		unitcd := cast.ToInt64(c.FormValue("unitcd"))

		ingamesn := c.FormValue("ingamesn")
		if len(ingamesn) == 0 && unitcd == UnitCD_NintendoWii {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid ingamesn",
				Err:     fmt.Errorf("wfc: protocol: legnth of ingamesn was 0"),
			})
		} else if len(ingamesn) != 0 {
			if goaway.IsProfane(ingamesn) {
				return NASReply(c, fiber.Map{
					"returncd": ReCD_ProfaneName,
				})
			}
		}

		devname := c.FormValue("devname")
		if len(devname) == 0 && unitcd == UnitCD_NintendoDS {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid ingamesn",
				Err:     fmt.Errorf("wfc: protocol: legnth of devname was 0"),
			})
		} else if len(devname) != 0 {
			if goaway.IsProfane(devname) {
				return NASReply(c, fiber.Map{
					"returncd": ReCD_ProfaneName,
				})
			}
		}

		return c.Next()
	}
}
