package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/qr/protocol"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
)

func qr2Delegate(conn *net.UdpPacket, logger *log.Logger) {

	qci := protocol.QR2CommandInfo{
		Log: logger.FactoryWithPostfix("QR2",
			fmt.Sprintf("<IP: %s>", conn.GetRemoteAddress()),
		),
	}

	qci.Log.Info("Handler", "bwaaa")
}
