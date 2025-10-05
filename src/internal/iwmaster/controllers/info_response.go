package controllers

import (
	"time"

	"buffersnow.com/spiritonline/internal/iwmaster/list"
	"buffersnow.com/spiritonline/internal/iwmaster/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
	"github.com/luxploit/red"
)

func handleInfoResponse(i *protocol.IWContext) error {
	lst, err := red.Locate[list.ServerList]()
	if err != nil {
		i.Log.Error("InfoResponse", "Failed to locate service: %v", err)
		return protocol.IWError_InvalidLocation
	}

	if len(i.CommandInfo.Data) != 1 {
		return protocol.IWError_InvalidCommand
	}

	kvs := gp.PickleMessage(i.CommandInfo.Data[0])
	game := gp.FindFromBundle("gamename", kvs).Value().String()
	challenge := gp.FindFromBundle("challenge", kvs).Value().String()

	err = lst.Access(game, challenge, func(s *list.Server) error {
		if s.Challenge != challenge {
			return protocol.IWError_ChallengeMismatch
		}

		proto, err := gp.FindFromBundle("protocol", kvs).Value().Integer()
		if err != nil {
			return protocol.IWError_InvalidProtocol
		}

		if s.Protocol != 0 && s.Protocol != proto {
			return protocol.IWError_InvalidProtocol
		}

		all_players, err := gp.FindFromBundle("clients", kvs).Value().Integer()
		if err != nil {
			return protocol.IWError_ValidationError
		}

		bot_players, err := gp.FindFromBundle("bots", kvs).Value().Integer()
		if err != nil {
			return protocol.IWError_ValidationError
		}

		hostname := gp.FindFromBundle("hostname", kvs).Value().String()
		if len(hostname) == 0 {
			return protocol.IWError_ValidationError
		}

		max_players, err := gp.FindFromBundle("sv_maxplayers", kvs).Value().Integer()
		if err != nil {
			return protocol.IWError_ValidationError
		}

		s.Protocol = proto
		s.State = list.ServerState_Idle
		s.Players = all_players - bot_players
		s.MaxPlayers = max_players
		s.Name = hostname
		s.LastPing = time.Now()

		i.Log.Info("Server List", "Successfully registered server %s for game %s", hostname, game)
		return nil
	})

	return err
}
