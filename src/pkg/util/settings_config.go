package util

import (
	"os"

	"buffersnow.com/spiritonline/pkg/version"
	"gopkg.in/yaml.v3"
)

type MySQLConfiguration struct {
	Host     string
	Port     int
	Username string
	Password string
	Database string
}

type SpiritConfiguration struct {
	Token        string `yaml:"server-token"`
	SkeletonHost string `yaml:"skeleton-host"`
	SkeletonPort int    `yaml:"skeleton-port"`
	RegionTag    string `yaml:"region-tag"`
}

type CommonConfiguration struct {
	Spirit SpiritConfiguration
	MySQL  MySQLConfiguration
}

func (settings utilSettings) loadConfig(config any) {

	Log.Info("Config", "Loading configuration file...")

	yamlFile, err := os.ReadFile(*settings.ConfigFolder + version.GetService() + ".yaml")
	if err != nil {
		Log.Panic("Config", "Failed to read config file!", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		Log.Panic("Config", "Failed to unmarshal config file!", err)
	}
}
