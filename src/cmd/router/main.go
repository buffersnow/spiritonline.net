package main

import (
	"buffersnow.com/spiritonline/internal/router"
	"buffersnow.com/spiritonline/pkg/di"
	"buffersnow.com/spiritonline/pkg/lifecycle"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/version"
)

func main() {
	version.PrintBuildInfo()
	app := di.New()

	app.Use(
		di.Provide(settings.New(&router.Config)),
		di.Provide(log.New),
		di.Invoke(lifecycle.AwaitInterrupt),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
