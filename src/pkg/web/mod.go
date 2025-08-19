package web

import (
	"errors"
	"fmt"
	"html/template"

	"buffersnow.com/spiritonline/pkg/log"
	"github.com/labstack/echo/v4"
	"github.com/luxploit/red"
)

type HttpUtils struct{}

func New() (*HttpUtils, error) {
	return &HttpUtils{}, nil
}

func (h HttpUtils) NewEcho(prefix string) (*echo.Echo, error) {
	tmpl, err := template.ParseGlob("public/*.html")
	if err != nil {
		return nil, fmt.Errorf("web: %w", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Renderer = &TemplateRenderer{Templates: tmpl}

	e.HTTPErrorHandler = ErrorHandler

	e.Use(XPoweredBy)
	e.Use(RequestLogging(prefix))

	e.RouteNotFound("/*", func(c echo.Context) error {
		return NotFoundError(&Details{
			Message: "seems like you took a wrong turn",
			Err:     errors.New("invalid resource"),
		})
	})

	return e, nil
}

func (h HttpUtils) StartEcho(e *echo.Echo, port int) (outerr error) {

	log, err := red.Locate[log.Logger]()
	if err != nil {
		return fmt.Errorf("web: %w", err)
	}

	// Echo doesn't always return an error so i'd rather have it catch the panic here,
	// and since we panic on error anyways atleast we know where it crashed (roughly)
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("web: echo: %v", r)
		}
	}()

	log.Info("HTTP Listener", "Listening on 0.0.0.0:%d", port)
	if err = e.Start(fmt.Sprintf(":%d", port)); err != nil {
		return fmt.Errorf("web: echo: %w", err)
	}

	return nil
}
