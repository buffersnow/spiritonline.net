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
)

func FieldsDecoder(sec *security.Security) echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if c.Request().Header.Get(echo.HeaderContentType) != echo.MIMEApplicationForm {
				return next(c)
			}

			if c.Request().Method != http.MethodPost {
				return next(c)
			}

			bodyBytes, err := io.ReadAll(c.Request().Body)
			if err != nil {
				return c.JSON(400, web.HttpError{
					Message: "bad body read",
					Error:   fmt.Sprintf("wfcnas: protocol: %v", err.Error()),
				})
			}

			err = c.Request().Body.Close()
			if err != nil {
				return c.JSON(400, web.HttpError{
					Message: "bad body close",
					Error:   fmt.Sprintf("wfcnas: protocol: %v", err.Error()),
				})
			}

			formVals, err := url.ParseQuery(string(bodyBytes))
			if err != nil {
				return c.JSON(400, web.HttpError{
					Message: "invalid form encoding",
					Error:   fmt.Sprintf("wfcnas: protocol: %v", err.Error()),
				})
			}

			for key, vals := range formVals {
				for i, v := range vals {
					decoded, err := sec.Encoding.DecodeB64_Wii([]byte(v))
					if err != nil {
						return c.JSON(400, web.HttpError{
							Message: "invalid base64 for field",
							Error:   fmt.Sprintf("wfcnas: protocol: %v", err.Error()),
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
}
