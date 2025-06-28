package settings

import "buffersnow.com/spiritonline/pkg/util"

type runtimeOptions struct {
	DisableDB   bool
	DBMigration bool
	LogArchival bool
	EnableDebug bool
	CertsFolder string
}

type Options struct {
	Runtime runtimeOptions

	MySQL struct {
		Host     string `env:"HOST,required"`
		Port     int    `env:"PORT" envDefault:"3306"`
		Username string `env:"USERNAME,required"`
		Password string `env:"PASSWORD,required"`
		Database string `env:"DATABASE" envDefault:"spiritonline"`
	} `envPrefix:"MYSQL_"`

	Spirit struct {
	} `envPrefix:"SPIRIT_"`

	Router struct {
	} `envPrefix:"ROUTER_"`
}

func New() (*Options, error) {
	settings := &Options{}

	tasks := []func() error{
		settings.loadEnv, // needs to be loaded first to avoid overriding flags
		settings.loadFlags,
	}

	return settings, util.Batch(tasks)
}
