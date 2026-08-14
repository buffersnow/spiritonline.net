package web

import (
	"encoding/xml"
	"errors"
	"fmt"

	"github.com/gofiber/fiber/v2"
	"github.com/luxploit/red"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/version"
)

type Details struct {
	XMLname xml.Name `xml:"PBXError"`
	Version string   `xml:"PBXVersion"`
	Code    int      `xml:"StatusCode"`
	Message string   `xml:"ErrMessage"`
	Err     error    `xml:"-"`
	Context Context  `xml:"Context"`
}

type baseTransmittableError struct {
	XMLName xml.Name `xml:"PBXError"`
	Version string   `xml:"PBXVersion"`
	Code    int      `xml:"StatusCode"`
	Message string   `xml:"ErrMessage"`
	Error   string   `xml:"SysError"`
	Context Context  `xml:"Context"`
}

type renderableError struct {
	Version string `xml:"PBXVersion"`
	Code    int    `xml:"StatusCode"`
	Message string `xml:"ErrMessage"`
	Error   string `xml:"SysError"`
	Context string `xml:"Context"`
}

func (e *Details) Error() string {
	return e.Err.Error()
}

func (e *Details) Unwrap() error {
	return e.Err
}

func renderError(c *fiber.Ctx, details *Details) error {
	if len(details.Version) <= 0 {
		details.Version = "??Unknown??"
	}

	//~ This is fucking stupid
	// if c.Get("X-Machine-Readable", "false") != "false" {
	return c.Status(details.Code).XML(details)
	// }

	//@ TODO: we dont yet support rendering html errors so always give out XML
	// return c.Status(details.Code).Render("error", renderableError{
	// 	Version: details.Version,
	// 	Code:    details.Code,
	// 	Message: details.Message,
	// 	Error:   details.Err.Error(),
	// 	Context: details.ContextXML(),
	// })
}

func ErrorHandler(c *fiber.Ctx, inerr error) error {
	build, err := red.Locate[version.BuildTag]()
	if err != nil {
		return renderError(c, &Details{
			Code:    fiber.StatusInternalServerError,
			Message: "bad service location attempt",
			Err:     fmt.Errorf("web: %w", err),
		})
	}

	logger, err := red.Locate[log.Logger]()
	if err != nil {
		return renderError(c, &Details{
			Code:    fiber.StatusInternalServerError,
			Message: "bad service location attempt",
			Err:     fmt.Errorf("web: %w", err),
			Version: build.GetFullTag(),
		})
	}

	var be *Details
	if errors.As(inerr, &be) {
		be.Version = build.GetFullTag()

		err := renderError(c, be)
		if err == nil {
			logger.Error("HTTP Error Handler",
				"<IP: %s> <Status: %d> Message: %s, Error: %s",
				c.IP(), be.Code, be.Message, be.Err.Error(),
			)
		} else {
			logger.Error("HTTP Error Handler",
				"<IP: %s> <Status: %d> Rendering Error: %s",
				c.IP(), c.Response().StatusCode(), err.Error(),
			)
		}

		return err
	}

	//& Fallback for unexpected errors
	logger.Error("HTTP Error Handler", "<IP: %s> <Status: 418> Unexpected Fallback! Error: %s", c.IP(), inerr.Error())
	return renderError(c, &Details{
		Code:    fiber.StatusTeapot,
		Message: "how did we get here?",
		Err:     inerr,
		Version: build.GetFullTag(),
		Context: Context{
			"RootCause": "this is a teapot, it can not brew coffee!",
		},
	})
}
