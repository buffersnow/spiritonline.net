package main

import (
	"buffersnow.com/spiritonline/pkg/lifecycle"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/version"

	"github.com/luxploit/red"
)

func main() {
	app := red.New()

	app.Use(
		red.Provide(version.New),
		red.Provide(settings.New),
		red.Provide(log.New),
		red.Invoke(lifecycle.Await),
	)

	if err := app.Run(); err != nil {
		panic(err)
	}
}
