package protocol

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"buffersnow.com/spiritonline/pkg/security"
	"buffersnow.com/spiritonline/pkg/web"
	"github.com/labstack/echo/v4"
	"github.com/luxploit/red"
)

func FieldsDecoder(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		if c.Request().Header.Get(echo.HeaderContentType) != echo.MIMEApplicationForm {
			return next(c)
		}

		if c.Request().Method != http.MethodPost {
			return next(c)
		}

		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return web.InternalServerError(&web.Details{
				Message: "bad body read",
				Err:     fmt.Errorf("wfcnas: protocol: %w", err),
			})
		}

		err = c.Request().Body.Close()
		if err != nil {
			return web.InternalServerError(&web.Details{
				Message: "bad body close",
				Err:     fmt.Errorf("wfcnas: protocol: %w", err),
			})
		}

		formVals, err := url.ParseQuery(string(bodyBytes))
		if err != nil {
			return web.BadRequestError(&web.Details{
				Message: "invalid form encoding",
				Err:     fmt.Errorf("wfcnas: protocol: %w", err),
			})
		}

		sec, err := red.Locate[security.Security]()
		if err != nil {
			return web.InternalServerError(&web.Details{
				Message: "bad service location",
				Err:     fmt.Errorf("wfcnas: protocol: %w", err),
			})
		}

		for key, vals := range formVals {
			for i, v := range vals {
				decoded, err := sec.Encoding.DecodeB64_Wii([]byte(v))
				if err != nil {
					return web.BadRequestError(&web.Details{
						Message: "invalid base64 for field",
						Err:     fmt.Errorf("wfcnas: protocol: %w", err),
					})

				}
				formVals[key][i] = string(decoded)
				println(key, i, string(decoded))
			}
		}

		// Rebuild the body with decoded values
		decodedBody := formVals.Encode()
		c.Request().Body = io.NopCloser(strings.NewReader(decodedBody))
		c.Request().ContentLength = int64(len(decodedBody))

		return nil
	}
}
