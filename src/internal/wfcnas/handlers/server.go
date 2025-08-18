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

	e, err := web.NewEcho("Nintendo WFC NAS")
	if err != nil {
		return fmt.Errorf("wfcnas: %w", err)
	}

	//@ TODO: Move conntest.nintendowifi.net to a seperate server, i know thats silly but still

	e.Use(protocol.XOrganization) // if this isn't here, conntest fails

	// ConnTest asks for /
	e.GET("/", controllers.NasTest)
	e.GET("/nastest.jsp", controllers.NasTest)

	g := e.Group("/", protocol.FieldsDecoder)
	{
		g.POST("/ac", controllers.Account)
		g.POST("/pr", controllers.Product)
		g.POST("/download", controllers.Download)
	}

	err = web.StartEcho(e, opt.Service.Ports["wfcnas"])
	if err != nil {
		return fmt.Errorf("wfcnas: %w", err)
	}

	return nil
}
