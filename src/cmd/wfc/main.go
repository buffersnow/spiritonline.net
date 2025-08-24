package main

import (
	"buffersnow.com/spiritonline/internal/wfc/handlers"
	"buffersnow.com/spiritonline/pkg/lifecycle"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/security"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/version"
	"buffersnow.com/spiritonline/pkg/web"

	"github.com/luxploit/red"
)

func main() {
	app := red.New()

	app.Use(
		red.Provide(version.New),
		red.Provide(settings.New),
		red.Provide(log.New),
		red.Provide(security.New),
		red.Provide(web.New),
		red.Invoke(handlers.ListenNASWii),
		red.Invoke(handlers.ListenNASDS),
		red.Invoke(handlers.ListenConntest),
		red.Invoke(handlers.ListenDls1),
		red.Invoke(lifecycle.Await),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
