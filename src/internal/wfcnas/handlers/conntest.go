package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfcnas/controllers"
	"buffersnow.com/spiritonline/internal/wfcnas/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenConntest(web *web.HttpUtils, opt *settings.Options, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: conntest: %w", err)
	}

	app.Use(protocol.XOrganization(), protocol.FieldsDecoder())

	app.Get("/", controllers.Test)
	app.Get("/nastest.jsp", controllers.Test) // WiiLink puts this here so god knows

	err = web.StartFiber(app, opt.Service.Ports["wfctest"])
	if err != nil {
		return fmt.Errorf("wfc: conntest: %w", err)
	}

	return nil
}
