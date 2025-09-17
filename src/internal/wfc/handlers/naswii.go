package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfc/controllers"
	"buffersnow.com/spiritonline/internal/wfc/protocol"
	"buffersnow.com/spiritonline/internal/wfc/repositories"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenNASWii(web *web.HttpUtils, repo *repositories.WFCRepo, opt *settings.Options, logger *log.Logger) error {

	app, err := web.NewFiber()
	if err != nil {
		return fmt.Errorf("wfc: nas-wii: %w", err)
	}

	app.Use(
		protocol.XOrganization(),
		protocol.FieldsDecoder(),
		protocol.ValidateRequest(),
		protocol.RequestFixup(),
		protocol.ProfanityFilter(),
		protocol.MarioKartOnly(),
	)

	app.Post("/ac", controllers.AccountWii)
	app.Post("/pr", controllers.Profanity)

	err = web.StartFiber(app, opt.Service.Ports["wfcwii"])
	if err != nil {
		return fmt.Errorf("wfc: nas-wii: %w", err)
	}

	return nil
}
