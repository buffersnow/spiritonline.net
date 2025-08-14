package web

import (
	"fmt"
	"html/template"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/version"
	"github.com/labstack/echo/v4"
)

type HttpError struct {
	Message string `json:"message"`
	Error   string `json:"error"`
}

type HttpUtils struct {
	bld *version.BuildTag
	log *log.Logger
}

func New(bld *version.BuildTag, log *log.Logger) (*HttpUtils, error) {
	return &HttpUtils{bld: bld, log: log}, nil
}

func (h HttpUtils) NewEcho(prefix string) (*echo.Echo, error) {
	tmpl, err := template.ParseGlob("public/*.html")
	if err != nil {
		return nil, fmt.Errorf("wfcnas: %w", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Renderer = &TemplateRenderer{Templates: tmpl}

	e.Use(RequestLogging(prefix, h.log))
	e.Use(XPoweredBy(h.bld))

	return e, nil
}

func (h HttpUtils) StartEcho(e *echo.Echo, port int) (outerr error) {

	// Echo doesn't always return an error so i'd rather have it catch the panic here,
	// and since we panic on error anyways atleast we know where it crashed (roughly)
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("echo: %v", r)
		}
	}()

	err := e.Start(fmt.Sprintf(":%d", port))
	if err != nil {
		return fmt.Errorf("echo: %w", err)
	}

	return nil
}
