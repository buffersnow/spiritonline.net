package handlers

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
	"buffersnow.com/spiritonline/pkg/settings"

	"buffersnow.com/spiritonline/internal/iwmaster/controllers"
	"buffersnow.com/spiritonline/internal/iwmaster/list"
)

func ListenService(opt *settings.Options, log *log.Logger, net *net.NetUtils, lst *list.ServerList) error {

	srv, err := net.CreateUdpListener(opt.Service.Ports["iwmaster"])
	if err != nil {
		return fmt.Errorf("qr: iw: %w", err)
	}

	go controllers.HandleIWMasterQueryServerInfo(lst)

	for {
		udp, err := srv.ReadBytes()
		if err != nil {
			log.Error("IWMaster Listener", "ReadBytes() failed: %v", err)
			continue
		}

		go controllers.HandleIWMasterIncoming(udp, log)
	}
}
