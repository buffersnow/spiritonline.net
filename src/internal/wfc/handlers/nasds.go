package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/controllers"
	"buffersnow.com/spiritonline/internal/wfc/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenNASDS(web *web.HttpUtils, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: nas-ds: %w", err)
	}

	app.Use(
		protocol.XOrganization(),
		protocol.FieldsDecoder(),
		protocol.ValidateRequest(),
		protocol.RequestFixup(),
		protocol.ProfanityFilter(),
	)

	app.Post("/ac", controllers.AccountDS)
	app.Post("/pr", controllers.Profanity)

	err = web.StartFiber(app, 5603) //@ TODO: FIX ME!
	if err != nil {
		return fmt.Errorf("wfc: nas-ds: %w", err)
	}

	return nil
}
