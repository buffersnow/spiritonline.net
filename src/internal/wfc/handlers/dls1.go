package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/controllers"
	"buffersnow.com/spiritonline/internal/wfc/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenDls1(web *web.HttpUtils, opt *settings.Options, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: dls1: %w", err)
	}

	app.Use(protocol.XOrganization(), protocol.FieldsDecoder())

	app.Get("/", controllers.Test)
	app.Get("/nastest.jsp", controllers.Test) // WiiLink puts this here so god knows

	err = web.StartFiber(app, opt.Service.Ports["wfcdls1"])
	if err != nil {
		return fmt.Errorf("wfc: dls1: %w", err)
	}

	return nil
}
