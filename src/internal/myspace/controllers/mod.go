package controllers

import (
	"buffersnow.com/spiritonline/internal/myspace/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
)

// (ctx *protocol.MySpaceContext, gci gp.GameSpyCommandInfo) error
type MSIMHandlerFunc = func(*protocol.MySpaceContext, gp.GameSpyCommandInfo) error

var command_routes = map[string]MSIMHandlerFunc{}

var callback_routes = map[string]MSIMHandlerFunc{}
