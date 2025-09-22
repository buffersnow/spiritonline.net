package controllers

import (
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

func HandleIncoming(g *protocol.GamespyContext, stream string) {

	data := strings.SplitSeq(stream, "final\\")
	for command := range data {
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

	// subCmdId := 0
	// if commandPair.Length() != 0 {
	// 	i64_subCmdId, err := commandPair.Value().Integer()
	// 	if err != nil {

	// 		return
	// 	}

	// }

}
