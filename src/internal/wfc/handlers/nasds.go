package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/controllers"
	"buffersnow.com/spiritonline/internal/wfc/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenNASDS(web *web.HttpUtils, opt *settings.Options, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: nas-ds: %w", err)
	}

	app.Use(protocol.XOrganization(), protocol.FieldsDecoder())

	app.Post("/ac", controllers.AccountDS)
	app.Post("/pr", controllers.Profanity)

	err = web.StartFiber(app, opt.Service.Ports["wfcds"])
	if err != nil {
		return fmt.Errorf("wfc: nas-ds: %w", err)
	}

	return nil
}
