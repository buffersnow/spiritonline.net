package settings

import (
	"fmt"
	"os"
	"reflect"
	"slices"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"

	"buffersnow.com/spiritonline/pkg/version"
)

var fs *pflag.FlagSet

func (o *Options) loadArgs(ver *version.BuildTag) error {

	fs = pflag.NewFlagSet("", pflag.ContinueOnError)
	fs.Usage = func() {
		fmt.Printf("Usage: ./%s <... arguments>\n", ver.GetService())
		o.helpText(reflect.ValueOf(o), "", "")
		fmt.Printf("\nThis help menu can be brought up again by running \"./%s help\"\n", ver.GetService())
		os.Exit(0)
	}

	if err := o.parseArgs(reflect.ValueOf(o)); err != nil {
		return fmt.Errorf("settings: parser: %w", err)
	}

	err := fs.Parse(os.Args[1:])
	if err != nil && err != pflag.ErrHelp {
		return fmt.Errorf("settings: %w", err)
	} else if err == pflag.ErrHelp {
		fs.Usage()
	}

	if slices.Contains(fs.Args(), "help") {
		fs.Usage()
	}

	o.Runtime.LogArchival = !o.Runtime.LogArchival

	return nil
}

func (o *Options) loadEnv(ver *version.BuildTag) error {

	err := godotenv.Load(fmt.Sprintf(".env.%s", ver.GetService()))
	if err != nil {
		return fmt.Errorf("settings: env: os: %w", err) // the only error thrown is from os.Open in godotenv.go L#207
	}

	var cfg Options
	err = env.Parse(&cfg)
	if err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	o.MySQL = cfg.MySQL
	o.Spirit = cfg.Spirit
	o.Service = cfg.Service

	return nil
}
