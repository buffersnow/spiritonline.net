package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
)

func ListenGPSP(net *net.NetUtils, log *log.Logger) error {
	srv, err := net.CreateTcpListener(29901) //@ TODO: FIX ME!
	if err != nil {
		return fmt.Errorf("gsp: gpsp: %w", err)
	}

	for {
		cli, err := srv.Accept()
		if err != nil {
			log.Error("GPSP Listener", "Accept() failed: %v", err)
			continue
		}

		go gpspDelegate(cli, log)
	}
}

func gpspDelegate(conn *net.TcpConnection, logger *log.Logger) {

	ctx := &protocol.GamespyContext{
		Connection: conn,
		Log: logger.FactoryWithPostfix("GPSP",
			fmt.Sprintf("<IP: %s>", conn.GetRemoteAddress()),
		),
	}

	ctx.Log.Event("Client", "Client awaiting authentication!")

	defer func() {
		ctx.Log.Event("Client", "Client exited!")
		ctx.Connection.Close()
	}()

	for {
		_, err := conn.ReadText()
		if err != nil {
			ctx.Log.Debug(log.DEBUG_TRAFFIC, "Server", "Traffic read error debug: %v", err)
			break
		}

	}

}
