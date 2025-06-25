package router

import (
	"buffersnow.com/spiritonline/pkg/settings"
)

type RouterConfiguration struct {
	Common settings.CommonConfiguration
}

var Config RouterConfiguration
