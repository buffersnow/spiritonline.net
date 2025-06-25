package settings

import (
	"time"
)

type Options struct {
	TestNoDb           *bool
	MigrateDB          *bool
	NoLogArchival      *bool
	LogFileName        *string
	CompressionJobTime *time.Duration
	ShowServerDebug    *bool
	ConfigFolder       *string
	CertsFolder        *string
	Standalone         *bool
	ReconnectDelay     *int
}

func New(config any) func() (*Options, error) {
	return func() (*Options, error) {
		settings := &Options{}

		settings.scanFlags()
		err := settings.loadConfig(config)

		return settings, err
	}
}
