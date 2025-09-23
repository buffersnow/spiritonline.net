package controllers

import (
	"errors"
	"strings"

	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/log"
)

type GPHandlerFunc = func(*protocol.GamespyContext, gp.GameSpyCommandInfo) error

var gpcm_routes = map[string]GPHandlerFunc{
	protocol.GPCommand_KeepAlive:         func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCommand_Error:             func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_Login:           func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_Logout:          func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_UpdateProfile:   func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_UpdateStatus:    func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_AddBuddy:        func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_DeleteBuddy:     func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_AuthorizeFriend: func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_BuddyMessage:    func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCMCommand_GetProfile:      func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
}

var gpsp_routes = map[string]GPHandlerFunc{
	protocol.GPCommand_KeepAlive:    func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPCommand_Error:        func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPSPCommand_OthersList: func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
	protocol.GPSPCommand_Search:     func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) error { return nil },
}

func StartGPCMAuth(g *protocol.GamespyContext) {
	g.SendRaw([]gp.GameSpyKV{
		gp.Message.Integer("lc", 1),
		gp.Message.String("challenge", g.GPCM.Challenge),
		gp.Message.Integer("id", 1),
	})
}

func HandleIncoming(g *protocol.GamespyContext, stream string) error {

	data := strings.SplitSeq(stream, "final\\")
	for command := range data {
		kvs := gp.PickleMessage(command)

		err := handleClientCommands(g, kvs)
		if err == nil {
			continue
		}

		var gpErr *gp.GameSpyError
		if errors.As(err, &gpErr) {
			g.Error(*gpErr)
			return err
		}

		return err
	}

	return nil
}

func handleClientCommands(g *protocol.GamespyContext, kvs []gp.GameSpyKV) error {

	if len(kvs) == 0 {
		return errors.New("gsp: parser: no kv's found")
	}

	commandPair := kvs[0] //& should always be the command
	g.Log.Trace(log.DEBUG_SERVICE, "Parser", "Processing command %s", commandPair.Key())

	subCmdId := 0
	if commandPair.Length() != 0 {
		id, err := commandPair.Value().Integer()
		if err != nil {
			return protocol.GPError_Parse
		}

		subCmdId = id
		g.Log.Trace(log.DEBUG_SERVICE, "Parser", "Sub-Command Id: %d", subCmdId)
	}

	gci := gp.GameSpyCommandInfo{
		Command:    commandPair.Key(),
		SubCommand: subCmdId,
		Data:       kvs[1:],
	}

	//& this means its a GPCM server message, and not GPSP
	if len(g.GPCM.Challenge) != 0 {
		h, ok := gpcm_routes[gci.Command]
		if !ok {
			return protocol.GPError_UnknownCommand
		}

		return h(g, gci)
	}

	h, ok := gpsp_routes[gci.Command]
	if !ok {
		return protocol.GPError_UnknownCommand
	}

	return h(g, gci)
}
