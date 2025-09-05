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
	SessionKey   uint
	BuildNumber  uint
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
	Body           iter.Seq[gp.GameSpyKV]
}

func (cbInfo MySpaceCallbackInfo) Find(key string) gp.GameSpyKV {
	return gp.FindFromBundle(key, cbInfo.Body)
}
