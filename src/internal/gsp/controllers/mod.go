package controllers

import (
	"slices"
	"strings"

	"buffersnow.com/spiritonline/internal/gsp/protocol"
	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/log"
)

var gpcm_routes = map[string]func(*protocol.GamespyContext, gp.GamespyCommandInfo){
	protocol.GPCMCommand_KeepAlive:       func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_Login:           func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_Logout:          func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_UpdateProfile:   func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_UpdateStatus:    func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_AddBuddy:        func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_DeleteBuddy:     func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_AuthorizeFriend: func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_BuddyMessage:    func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPCMCommand_GetProfile:      func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
}

var gpsp_routes = map[string]func(*protocol.GamespyContext, gp.GamespyCommandInfo){
	protocol.GPSPCommand_KeepAlive:  func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPSPCommand_OthersList: func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
	protocol.GPSPCommand_Search:     func(g *protocol.GamespyContext, gci gp.GamespyCommandInfo) {},
}

func HandleIncoming(g *protocol.GamespyContext, stream string) {

	data := strings.Split(stream, "final\\")
	for _, command := range data {
		kvs := gp.PickleMessage(command)
		handleClientCommands(g, slices.Collect(kvs))
	}
}

func handleClientCommands(g *protocol.GamespyContext, kvs []gp.GameSpyKV) {

	if len(kvs) == 0 {
		return
	}

	commandPair := kvs[0] //& should always be the command
	g.Log.Trace(log.DEBUG_SERVICE, "Parser", "Processing command %s", commandPair.Key())

}
