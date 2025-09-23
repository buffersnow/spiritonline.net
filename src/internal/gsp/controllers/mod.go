package controllers

import (
	"errors"
	"fmt"
	"slices"
	"strings"

	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/log"
)

var gpcm_routes = map[string]func(*protocol.GamespyContext, gp.GameSpyCommandInfo){
	protocol.GPCMCommand_KeepAlive:       func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_Login:           func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_Logout:          func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_UpdateProfile:   func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_UpdateStatus:    func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_AddBuddy:        func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_DeleteBuddy:     func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_AuthorizeFriend: func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_BuddyMessage:    func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPCMCommand_GetProfile:      func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
}

var gpsp_routes = map[string]func(*protocol.GamespyContext, gp.GameSpyCommandInfo){
	protocol.GPSPCommand_KeepAlive:  func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPSPCommand_OthersList: func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
	protocol.GPSPCommand_Search:     func(g *protocol.GamespyContext, gci gp.GameSpyCommandInfo) {},
}

func StartGPCMAuth(g *protocol.GamespyContext) {
	g.Send([]gp.GameSpyKV{
		gp.Message.Integer("lc", 1),
		gp.Message.String("challenge", g.GPCM.Nonce),
		gp.Message.Integer("id", 1),
	})
}

func HandleIncoming(g *protocol.GamespyContext, stream string) error {

	data := strings.SplitSeq(stream, "final\\")
	for command := range data {
		kvs := gp.PickleMessage(command)

		err := handleClientCommands(g, slices.Collect(kvs))
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
			return fmt.Errorf("gsp: parser: %w", err)
		}

		subCmdId = id
	}

	g.Log.Trace(log.DEBUG_SERVICE, "Parser", "Sub-Command Id: %d", subCmdId)

	return nil
}
