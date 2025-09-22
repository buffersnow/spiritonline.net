package protocol

import (
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
