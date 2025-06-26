package main

import (
	"buffersnow.com/spiritonline/internal/router"
	"buffersnow.com/spiritonline/pkg/lifecycle"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/version"
	"github.com/luxploit/red"
)

func main() {
	version.PrintBuildInfo()
	app := red.New()

	app.Use(
		red.Provide(settings.New(&router.Config)),
		red.Provide(log.New),
		red.Invoke(lifecycle.AwaitInterrupt),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
