package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/controllers"
	"buffersnow.com/spiritonline/internal/wfc/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenConntest(web *web.HttpUtils, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: conntest: %w", err)
	}

	app.Use(protocol.XOrganization())

	app.Get("/", controllers.Test)
	app.Get("/nastest.jsp", controllers.Test) //? WiiLink puts this here so god knows

	err = web.StartFiber(app, 5601) //@ TODO: FIX ME!
	if err != nil {
		return fmt.Errorf("wfc: conntest: %w", err)
	}

	return nil
}
