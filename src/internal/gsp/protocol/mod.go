package protocol

import (
	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
)

type GamespyContext struct {
	Connection *net.TcpConnection
	GPCM       GPCMContext
	Log        log.LoggingFactory
}

type GPCMContext struct {
	Nonce      string
	SessionKey int
}

func StartGPCMAuth(g *GamespyContext) {
	g.Send([]gp.GameSpyKV{
		gp.Message.Integer("lc", 1),
		gp.Message.String("challenge", g.GPCM.Nonce),
		gp.Message.Integer("id", 1),
	})
}
