package handlers

import (
	"fmt"
	"math/rand/v2"

	"buffersnow.com/spiritonline/internal/gsp/controllers"
	"buffersnow.com/spiritonline/internal/gsp/protocol"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
	"buffersnow.com/spiritonline/pkg/settings"
	"buffersnow.com/spiritonline/pkg/util"
)

func ListenGPCM(opt *settings.Options, net *net.NetUtils, log *log.Logger) error {
	srv, err := net.CreateTcpListener(opt.Service.Ports["gpcm"])
	if err != nil {
		return fmt.Errorf("gsp: gpcm: %w", err)
	}

	for {
		cli, err := srv.Accept()
		if err != nil {
			log.Error("GPCM Listener", "Accept() failed: %v", err)
			continue
		}

		go gpcmDelegate(cli, log)
	}
}

func gpcmDelegate(conn *net.TcpConnection, logger *log.Logger) {

	ctx := &protocol.GamespyContext{
		Connection: conn,
		GPCM: protocol.GPCMContext{
			Challenge:  util.RandomString(0x40),
			SessionKey: rand.IntN(0xFFFFF),
		},
		Log: logger.FactoryWithPostfix("GPCM",
			fmt.Sprintf("<IP: %s>", conn.GetRemoteAddress()),
		),
	}

	ctx.Log.Event("Client", "Client awaiting authentication!")

	defer func() {
		ctx.Log.Event("Client", "Client exited!")
		ctx.Connection.Close()
	}()

	controllers.StartGPCMAuth(ctx)

	for {
		_, err := conn.ReadText()
		if err != nil {
			ctx.Log.Debug(log.DEBUG_TRAFFIC, "Server", "Traffic read error debug: %v", err)
			break
		}

	}

}
