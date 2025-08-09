package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
	"buffersnow.com/spiritonline/pkg/settings"
)

func ListenService(opt *settings.Options, log *log.Logger) error {

	srv, err := net.CreateTcpListener(opt.Service.ProtocolPort, log)
	if err != nil {
		return fmt.Errorf("nas: %w", err)
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

}
