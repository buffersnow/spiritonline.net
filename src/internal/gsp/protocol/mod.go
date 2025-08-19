package protocol

import (
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
)

type GamespyContext struct {
	Connection *net.TcpConnection
	Client     GamespyClientContext
	Log        log.LoggingFactory
}

type GamespyClientContext struct {
	Nonce      string
	SessionKey int
}
