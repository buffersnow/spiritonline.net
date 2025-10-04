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

func HandleIWMasterQueryServerInfo(lst *list.ServerList) {
	for {
		lst.Iterate(func(game string, s *list.Server) {
			curTime := time.Now()

			if s.State == list.ServerState_Idle && curTime.After(s.LastPing.Add(15*time.Minute)) ||
				s.State == list.ServerState_Looking && curTime.After(s.LastPing.Add(2*time.Minute)) {

				lst.Remove(game, s)
				return
			}

			if s.State != list.ServerState_Refreshing {
				return
			}

			lst.Lock(func() {
				s.LastPing = curTime
				s.State = list.ServerState_Looking

				s.Context.Send(protocol.IWCommandInfo{
					Command: protocol.IWCommand_GetInfo,
					Data: []string{
						s.Challenge,
					},
				})
			})
		})

		time.Sleep(100 * time.Millisecond)
	}
}
