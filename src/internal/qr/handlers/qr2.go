package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/internal/qr/protocol"
	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
)

func qr2Delegate(conn *net.UdpPacket, logger *log.Logger) {

	ctx := protocol.QR2Context{
		Log: logger.FactoryWithPostfix("QR2",
			fmt.Sprintf("<IP: %s>", conn.GetRemoteAddress()),
		),
	}

	ctx.Log.Info("Handler", "bwaaa")
}
