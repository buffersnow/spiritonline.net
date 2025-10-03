package controllers

import (
	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
)

func handleGP_Error(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error {
	return protocol.GPError_DisallowedCommand
}
