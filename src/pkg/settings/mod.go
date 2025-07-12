package settings

import (
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/version"
)

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
		AuthorizationToken string `env:"AUTH_TOKEN,required"`
		HeadunitHost       string `env:"HEADUNIT_HOST,required"`
		HeadunitPort       int    `env:"HEADUNIT_PORT" envDefault:"1390"`
		ServiceTag         string `env:"SERVICE_TAG" envDefault:"ww-global-unknown-1"`
	} `envPrefix:"SPIRIT_"`

	Service struct {
		ProtocolPort int             `env:"PROTOCOL_PORT,required"`
		HttpPort     int             `env:"HTTP_PORT" envDefault:"9999"`
		Features     map[string]bool `env:"FEATURES,required"`
	} `envPrefix:"SERVICE_"`
}

func New(ver *version.BuildTag) (*Options, error) {
	settings := &Options{}

	tasks := []func() error{
		settings.loadEnv(ver), // needs to be loaded first to avoid overriding flags
		settings.loadFlags,
	}

	return settings, util.Batch(tasks)
}
