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

		//& "GameSpy implemented this on the same port,
		//& so version is determined by the presence of
		//& a '\' as the first byte." - OpenSpy Core v2
		//$ https://github.com/openspy/openspy-core
		//~ omw to fucking shoot myself with these fucking
		//~ '\\' hauting me from myspace to iw8 lan (luxploit)
		if udp.Data[0] == '\\' /*legacy*/ {
			//% handlers/qr1.go
			go qr1Delegate(udp, log)
			continue
		}

		//% handlers/qr2.go
		go qr2Delegate(udp, log)
	}
}
