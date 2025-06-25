package settings

import (
	"fmt"
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

// @ TODO: Maybe make a standardized config instead?
func (settings Options) loadConfig(config any) error {

	yamlFile, err := os.ReadFile(*settings.ConfigFolder + "/" + version.GetService() + ".yaml")
	if err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	err = yaml.Unmarshal(yamlFile, config)
	if err != nil {
		return fmt.Errorf("settings: %w", err)
	}

	return nil
}
