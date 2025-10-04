package controllers

import (
	"time"

	"buffersnow.com/spiritonline/internal/iwmaster/list"
	"buffersnow.com/spiritonline/internal/iwmaster/protocol"
	"buffersnow.com/spiritonline/pkg/log"

	"github.com/luxploit/red"
)

func handleHeartbeat(i *protocol.IWContext) error {

	lst, err := red.Locate[list.ServerList]()
	if err != nil {
		i.Log.Error("Heartbeat", "Failed to locate service: %v", err)
		return protocol.IWError_InvalidLocation
	}

	if len(i.CommandInfo.Data) <= 1 {
		return protocol.IWError_InvalidCommand
	}

	game := i.CommandInfo.Data[0]
	challenge := i.CommandInfo.Data[1]

	err = lst.Access(game, challenge, func(s *list.Server) error {
		lst.Lock(func() {
			s.State = list.ServerState_Refreshing
			s.LastPing = time.Now()
		})

		return nil
	})

	if err != nil {
		lst.Add(game, &list.Server{
			State:     list.ServerState_Refreshing,
			Challenge: challenge,
			LastPing:  time.Now(),
			Context:   i,
		})
	}

	i.Log.Trace(log.DEBUG_SERVICE, "Heartbeat", "Recieved heartbeat for server game type %s", game)

	return nil
}
