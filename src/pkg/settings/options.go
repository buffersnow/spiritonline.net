package settings

import (
	"errors"
	"fmt"
	"os"
	"reflect"
	"slices"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"github.com/spf13/pflag"

	"buffersnow.com/spiritonline/pkg/util"
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

	if len(o.Service.Ports) == 0 {
		return errors.New("settings: no ports were defined in the configuration")
	}

	return nil
}

func (o *Options) loadEnv() error {

	if err := godotenv.Load(); err != nil {
		return fmt.Errorf("settings: env: os: %w", err) // the only error thrown is from os.Open in godotenv.go L#207
	}

	// cleanup service_ports and service_features
	tasks := []func() error{
		func() error { return util.CleanEnv("SERVICE_PORTS") },
		func() error { return util.CleanEnv("SERVICE_FEATURES") },
		func() error { return util.CleanEnv("SERVICE_PROXIES") },
	}

	if err := util.Batch(tasks); err != nil {
		return fmt.Errorf("settings: env: %w", err)
	}

	var cfg Options
	if err := env.Parse(&cfg); err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	o.MySQL = cfg.MySQL
	o.Spirit = cfg.Spirit
	o.Service = cfg.Service

	return nil
}
