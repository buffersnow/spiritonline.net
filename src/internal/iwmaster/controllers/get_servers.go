package controllers

import (
	"bytes"
	"encoding/binary"

	"buffersnow.com/spiritonline/internal/iwmaster/list"
	"buffersnow.com/spiritonline/internal/iwmaster/protocol"
	"github.com/luxploit/red"
	"github.com/spf13/cast"
)

func getFilter(data []string, filter string) bool {
	if len(data) <= 1 {
		if data[0] == filter {
			return true
		}
	}

	if len(data) <= 2 {
		if data[1] == filter {
			return true
		}
	}

	return false
}

func handleGetServers(i *protocol.IWContext) error {

	lst, err := red.Locate[list.ServerList]()
	if err != nil {
		i.Log.Error("Server List", "Failed to locate service: %v", err)
		return protocol.IWError_InvalidLocation
	}

	if len(i.CommandInfo.Data) < 2 {
		return protocol.IWError_InvalidCommand
	}

	game := i.CommandInfo.Data[0]
	proto, err := cast.ToIntE(i.CommandInfo.Data[1])

	if len(game) == 0 || err != nil {
		return protocol.IWError_InvalidCommand
	}

	showFull := getFilter(i.CommandInfo.Data, "full")
	showEmpty := getFilter(i.CommandInfo.Data, "empty")

	filteredSrvs := []*list.Server{}
	lst.IterateRead(func(g string, s *list.Server) {
		if game != g || s.Protocol != proto {
			return
		}

		filteredSrvs = append(filteredSrvs, s)
	})

	pktCount := 0
	buf := bytes.Buffer{}
	for idx := 0; idx < len(filteredSrvs); idx++ {
		if !showEmpty && filteredSrvs[idx].Players == 0 {
			continue
		}

		if !showFull && filteredSrvs[idx].MaxPlayers == filteredSrvs[idx].Players {
			continue
		}

		buf.WriteByte('\\')

		ip4 := filteredSrvs[idx].Context.Connection.Addr.IP.To4()
		buf.Write(ip4[:4])

		binary.Write(&buf, binary.LittleEndian, uint16(filteredSrvs[idx].Context.Connection.Addr.Port))

		if buf.Len() >= 1400 || idx+1 == len(filteredSrvs) {
			if idx+1 == len(filteredSrvs) {
				buf.WriteByte('\\')
				buf.WriteString("EOT")
				buf.WriteByte(0x00)
				buf.WriteByte(0x00)
				buf.WriteByte(0x00)
			}

			i.Send(protocol.IWCommandInfo{
				Command: protocol.IWCommand_ServersResponse,
				Data: []string{
					buf.String(),
				},
			})

			pktCount++
			buf = bytes.Buffer{}
		}

	}

	i.Log.Info("Server List", "Sent %d servers in %d parts for game %s", len(filteredSrvs), pktCount, game)

	return nil
}
