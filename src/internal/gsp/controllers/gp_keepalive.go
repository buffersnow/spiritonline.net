package controllers

import (
	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
)

func handleGP_KeepAlive(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error {
	return g.Send(gci) // basically just echo back the ka
}
