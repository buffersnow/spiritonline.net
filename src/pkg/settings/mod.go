package settings

import (
	"buffersnow.com/spiritonline/pkg/util"
	"buffersnow.com/spiritonline/pkg/version"
)

//& Runtime.LogArchive gets inverted after parsing

type Options struct {
	Runtime struct {
		LogArchival bool   `arg:"--no-archive" help:"Disables logfile daily archive and compressions"`
		CertsFolder string `arg:"--certs-folder" default:"certs" help:"ECDSA public and private key directory"`
	}

	Development struct {
		EnableDev     bool `arg:"--devel" help:"Enable development mode"`
		DisableDB     bool `arg:"--no-database" help:"Run without a database for development"`
		EnableVerbose bool `arg:"-v,--verbose" help:"Show trace prints for extended logging"`
		EnableDebug   bool `arg:"-g,--vdebug" help:"Show all debug/developer log prints"`
	}

	DeploymentOverride struct {
		EnableOverride     bool   `arg:"--override" help:"Enable deployment override for BufferSnow internal deployment scripts"`
		AuthorizationToken string `arg:"-t,--auth-token" default:"randomgarbage" help:"Spirit network auth token override"`
		ServiceTag         string `arg:"-i,--service-tag" default:"ww-global-unknown-1" help:"Spirit service service-tag/id override"`
	} `section:"Deployment Override"`

	MySQL struct {
		Host     string `env:"HOST,required" help:"MySQL database server IP/Hostname"`
		Port     int    `env:"PORT" envDefault:"3306" help:"MySQL database server port"`
		Username string `env:"USERNAME,required" help:"MySQL database account username"`
		Password string `env:"PASSWORD,required" help:"MySQL database account Password"`
		Database string `env:"DATABASE" envDefault:"spiritonline" help:"MySQL database for services"`
	} `envPrefix:"MYSQL_"`

	Spirit struct {
		AuthorizationToken string `env:"AUTH_TOKEN,required" help:"Microservice network auth token"`
		HeadunitHost       string `env:"HEADUNIT_HOST,required" help:"Microservice headunit/router server IP/Hostname"`
		HeadunitPort       int    `env:"HEADUNIT_PORT" envDefault:"1390" help:"Microservice Headunit server port"`
		ServiceTag         string `env:"SERVICE_TAG" envDefault:"ww-global-unknown-1" help:"Microservice identification service-tag/id"`
	} `envPrefix:"SPIRIT_" section:"Spirit/Microservice"`

	Service struct {
		Ports    map[string]int    `env:"PORTS,required" help:"Service-specific protocol/http/misc ports"`
		Features map[string]bool   `env:"FEATURES,required" help:"Service-specific feature configuration list"`
		Proxies  map[string]string `env:"HTTP_PROXIES" help:"HTTP-only services specific reverse proxy list"`
	} `envPrefix:"SERVICE_"`
}

func New(ver *version.BuildTag) (*Options, error) {
	options := &Options{}

	tasks := []func() error{
		//! needs to be loaded first to avoid overriding flags
		func() error { return options.loadEnv() },
		func() error { return options.loadArgs(ver) },
	}

	return options, util.Batch(tasks)
}
