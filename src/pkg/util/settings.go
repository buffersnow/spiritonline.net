package util

import (
	"sync"
	"time"
)

type utilSettings struct {
	TestNoDb           *bool
	MigrateDB          *bool
	NoLogCompression   *bool
	LogFileName        *string
	CompressionJobTime *time.Duration
	ShowServerDebug    *bool
	ConfigFolder       *string
	CertsFolder        *string
	Standalone         *bool
	ReconnectDelay     *int
}

var Settings *utilSettings = &utilSettings{}
var settingsWg = sync.WaitGroup{}

func (s utilSettings) Initialize(config any) {
	s.scanFlags()
	s.loadConfig(config)
	Settings = &s
	settingsWg.Done()
}
