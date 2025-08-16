package main

import (
	"buffersnow.com/spiritonline/pkg/lifecycle"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
	"buffersnow.com/spiritonline/pkg/security"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/version"

	"buffersnow.com/spiritonline/internal/gsp/handlers"

	"github.com/luxploit/red"
)

func main() {
	app := red.New()

	app.Use(
		red.Provide(version.New),
		red.Provide(settings.New),
		red.Provide(log.New),
		red.Provide(security.New),
		red.Provide(net.New),
		red.Invoke(handlers.ListenGPCM),
		red.Invoke(handlers.ListenGPSP),
		red.Invoke(lifecycle.Await),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
