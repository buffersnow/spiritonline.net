package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfcnas/controllers"
	"buffersnow.com/spiritonline/internal/wfcnas/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenNAS(web *web.HttpUtils, opt *settings.Options, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: nas: %w", err)
	}

	app.Use(protocol.XOrganization(), protocol.FieldsDecoder())

	app.Post("/ac", controllers.Account)
	app.Post("/pr", controllers.Profanity)

	err = web.StartFiber(app, opt.Service.Ports["wfcnas"])
	if err != nil {
		return fmt.Errorf("wfc: nas: %w", err)
	}

	return nil
}
