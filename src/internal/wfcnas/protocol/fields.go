package protocol

import (
	"fmt"
	"net/url"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/security"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"
)

func FieldsDecoder() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Get(fiber.HeaderContentType) != fiber.MIMEApplicationForm {
			return c.Next()
		}

		if c.Method() != fiber.MethodPost {
			return c.Next()
		}

		sec, err := red.Locate[security.Security]()
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad service location",
				Err:     fmt.Errorf("wfcnas: protocol: %w", err),
			})
		}

		logger, err := red.Locate[log.Logger]()
		if err != nil {
			return web.InternalServerError(c, &web.Details{
				Message: "bad service location",
				Err:     fmt.Errorf("wfcnas: protocol: %w", err),
			})
		}

		formVals, err := url.ParseQuery(string(c.Body()))
		if err != nil {
			return web.BadRequestError(c, &web.Details{
				Message: "invalid form encoding",
				Err:     fmt.Errorf("wfcnas: protocol: %w", err),
			})
		}

		c.Request().PostArgs().Reset() // clear args
		for key, vals := range formVals {
			for _, v := range vals {
				decoded, err := sec.Encoding.DecodeB64_Wii([]byte(v))
				if err != nil {
					return web.BadRequestError(c, &web.Details{
						Message: "invalid base64 for field",
						Err:     fmt.Errorf("wfcnas: protocol: %w", err),
					})

				}
				c.Request().PostArgs().Add(key, string(decoded))
				logger.Debug(log.DEBUG_SERVICE, "NAS Field Decoder", "%s=%s", key, string(decoded))
			}
		}

		return c.Next()
	}
}
