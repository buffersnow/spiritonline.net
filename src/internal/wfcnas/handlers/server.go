package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/wfcnas/controllers"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/web"
)

func ListenService(web *web.HttpUtils, opt *settings.Options, logger *log.Logger) (outerr error) {

	e, err := web.NewEcho("Nintendo WFC NAS")
	if err != nil {
		return fmt.Errorf("wfcnas: %w", err)
	}

	e.GET("/", controllers.Index)
	e.GET("/nastest.jsp", controllers.NasTest)

	// Echo doesn't always return an error so i'd rather have it catch the panic here,
	// and since we panic on error anyways atleast we know where it crashed (roughly)
	defer func() {
		if r := recover(); r != nil {
			outerr = fmt.Errorf("wfcnas: echo: %v", r)
		}
	}()

	logger.Info("HTTP Listener", "Listening on 0.0.0.0:%d", opt.Service.HttpPort)
	err = e.Start(fmt.Sprintf(":%d", opt.Service.HttpPort))
	if err != nil {
		return fmt.Errorf("wfcnas: echo: %w", err)
	}

	return nil
}
