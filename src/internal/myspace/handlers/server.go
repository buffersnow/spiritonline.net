package handlers

import (
	"fmt"
	"math/rand/v2"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/util"

	"buffersnow.com/spiritonline/internal/myspace/protocol"
)

func ListenMySpace(opt *settings.Options, log *log.Logger) error {

	srv, err := net.CreateTcpListener(opt.Service.ProtocolPort, log)
	if err != nil {
		return fmt.Errorf("msim: %w", err)
	}

	for {
		cli, err := srv.Accept()
		if err != nil {
			log.Error("Service Listener", "Accept() failed: %v", err)
			continue
		}

		go svcDelegate(cli, log)
	}
}

func svcDelegate(conn *net.TcpConnection, logger *log.Logger) {

	cli := protocol.MySpaceContext{
		Connection: conn,
		Client: protocol.MySpaceClientContext{
			Nonce:        util.RandomString(0x40),
			SessionKey:   rand.IntN(0xFFFFF),
			IgnoreTicket: true,
		},
		Log: logger.FactoryWithPostfix("MySpace",
			fmt.Sprintf("<IP: %s>", conn.GetRemoteAddress()),
		),
	}

	cli.Log.Info("Client", "Client awaiting authentication!")

	defer func() {
		cli.Log.Info("Client", "Client exited!")
		cli.Connection.Close()
	}()

	for {
		_, err := conn.ReadText()
		if err != nil {
			cli.Log.Debug(log.DEBUG_TRAFFIC, "Server", "Traffic read error debug: %v", err)
			break
		}

	}

}
