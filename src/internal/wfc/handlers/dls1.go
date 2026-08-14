package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/controllers"
	"buffersnow.com/spiritonline/internal/wfc/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenDls1(web *web.HttpUtils, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: dls1: %w", err)
	}

	app.Use(
		protocol.XOrganization(),
		protocol.FieldsDecoder(),
		protocol.ValidateRequest(),
		protocol.RequestFixup(),
	)

	app.Get("/download", controllers.Download)

	err = web.StartFiber(app, 5602) //@ TODO: FIX ME!
	if err != nil {
		return fmt.Errorf("wfc: dls1: %w", err)
	}

	return nil
}
