package handlers

import (
	"fmt"
	"html/template"

	"buffersnow.com/spiritonline/internal/wfcnas/controllers"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"

	"github.com/labstack/echo/v4"
)

func ListenService(opt *settings.Options, logger *log.Logger) error {
	tmpl, err := template.ParseGlob("public/*.html")
	if err != nil {
		return fmt.Errorf("wfcnas: %w", err)
	}

	e := echo.New()
	e.HideBanner = true
	e.HidePort = true
	e.Renderer = &web.TemplateRenderer{Templates: tmpl}

	e.Use(web.RequestLogging("Nintendo WFC NAS", logger))

	e.GET("/", controllers.Index)
	e.GET("/nastest.jsp", controllers.NasTest)

	logger.Info("TCP Listener", "Listening on 0.0.0.0:%d", opt.Service.HttpPort)
	err = e.Start(fmt.Sprintf(":%d", opt.Service.HttpPort))
	if err != nil {
		return fmt.Errorf("wfcnas: echo: %w", err)
	}

	return nil
}
