package protocol

import (
	"iter"

	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
)

type MySpaceContext struct {
	Connection *net.TcpConnection
	Client     MySpaceClientContext
	Profile    MySpaceProfileContext
	Log        log.LoggingFactory
}

type MySpaceClientContext struct {
	Nonce        string
	SessionKey   int
	BuildNumber  int
	IgnoreTicket bool
}

type MySpaceProfileContext struct {
	ImageData []byte
}

type MySpaceCallbackInfo struct {
	CommandType    int
	CommandFamily  int
	CommandSubcode int
	RequestId      int
	Body           iter.Seq[gp.GamespyKV]
}

func (cbInfo MySpaceCallbackInfo) Find(key string) gp.GamespyKV {
	return gp.FindInternal(key, cbInfo.Body)
}
