package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
	"buffersnow.com/spiritonline/pkg/settings"
)

func ListenService(opt *settings.Options, log *log.Logger, net *net.NetUtils) error {

	srv, err := net.CreateUdpListener(opt.Service.Ports["qr"])
	if err != nil {
		return fmt.Errorf("qr: %w", err)
	}

	for {
		udp, err := srv.ReadBytes()
		if err != nil {
			log.Error("QR Listener", "ReadBytes() failed: %v", err)
			continue
		}

		go qrDelegate(udp, log)
	}
}

func qrDelegate(conn *net.UdpPacket, logger *log.Logger) {

}
