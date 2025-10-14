package controllers

import (
	"errors"
	"fmt"
	"time"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"

	"buffersnow.com/spiritonline/internal/iwmaster/list"
	"buffersnow.com/spiritonline/internal/iwmaster/protocol"
)

type iwRemovalList struct {
	game   string
	server *list.Server
}

// (i *protocol.IWMasterContext) error
type IWHandlerFunc = func(*protocol.IWContext) error

var iw_routes = map[string]IWHandlerFunc{
	protocol.IWCommand_Heartbeat:    handleHeartbeat,
	protocol.IWCommand_InfoResponse: handleInfoResponse,
	protocol.IWCommand_GetServers:   handleGetServers,
	protocol.IWCommand_GetBots:      handleGetBots,
	protocol.IWCommand_GetCRM:       handleGetCRM,
}

func HandleIWMasterIncoming(conn *net.UdpPacket, logger *log.Logger) {

	ctx := &protocol.IWContext{
		Log: logger.FactoryWithPostfix("IWMaster",
			fmt.Sprintf("<IP: %s>", conn.GetRemoteAddress()),
		),
		Connection:  conn,
		CommandInfo: protocol.PickleMessage(conn.Data),
	}

	ctx.Log.Debug(log.DEBUG_SERVICE, "Parser", "PickleMessage returned: %+v", ctx.CommandInfo)
	ctx.Log.Event("Parser", "Processing command %s", ctx.CommandInfo.Command)

	h, ok := iw_routes[ctx.CommandInfo.Command]
	if !ok {
		ctx.Log.Error("Parser", "Unknown command: %s", ctx.CommandInfo.Command)
		return
	}

	if err := h(ctx); err != nil {
		var iwErr *protocol.IWError
		if errors.As(err, &iwErr) {
			ctx.Error(iwErr)
		}

		ctx.Log.Error("Parser", "An error has occurred: %v", err)
	}
}

func HandleIWMasterWatchdog(logger *log.Logger, lst *list.ServerList) {
	for {
		var removeList []iwRemovalList
		var changeList []*list.Server

		lst.Iterate(func(game string, s *list.Server) {
			curTime := time.Now()

			if (s.State == list.ServerState_Idle && curTime.After(s.LastPing.Add(16*time.Minute))) ||
				(s.State == list.ServerState_Looking && curTime.After(s.LastPing.Add(2*time.Minute))) {
				removeList = append(removeList, iwRemovalList{game: game, server: s})
				return
			}

			if s.State == list.ServerState_Refreshing {
				changeList = append(changeList, s)
			}
		})

		for _, r := range removeList {
			lst.Remove(r.game, r.server)
			logger.Action("IWMaster Watchdog", "<IP: %s> Removing server for game type %s for inactivity",
				r.server.Context.Connection.GetRemoteAddress(), r.game,
			)
		}

		for _, s := range changeList {
			lst.Lock(func() {
				s.LastPing = time.Now()
				s.State = list.ServerState_Looking
				s.Context.Send(protocol.IWCommandInfo{
					Command: protocol.IWCommand_GetInfo,
					Data:    []string{s.Challenge},
				})
			})

			logger.Event("IWMaster Watchdog", "<IP: %s> Pinging idle server",
				s.Context.Connection.GetRemoteAddress(),
			)
		}

		time.Sleep(2 * time.Minute)
	}
}
