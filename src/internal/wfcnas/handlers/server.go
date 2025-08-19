package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfcnas/controllers"
	"buffersnow.com/spiritonline/internal/wfcnas/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenService(web *web.HttpUtils, opt *settings.Options, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfcnas: %w", err)
	}

	//@ TODO: Move conntest.nintendowifi.net to a seperate server, i know thats silly but still

	app.Use(protocol.XOrganization()) // if this isn't here, conntest fails

	// ConnTest asks for /
	app.Get("/", controllers.NasTest)
	app.Get("/nastest.jsp", controllers.NasTest) // WiiLink puts this here so god knows

	g := app.Group("", protocol.FieldsDecoder())
	{
		g.Post("/ac", controllers.Account)
		g.Post("/pr", controllers.Product)
		g.Post("/download", controllers.Download)
	}

	err = web.StartFiber(app, opt.Service.Ports["wfcnas"])
	if err != nil {
		return fmt.Errorf("wfcnas: %w", err)
	}

	return nil
}
