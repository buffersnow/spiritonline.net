package web

import (
	"fmt"
	"html/template"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/version"
	"github.com/labstack/echo/v4"
)

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
