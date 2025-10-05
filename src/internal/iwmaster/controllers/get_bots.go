package controllers

import (
	"os"

	"buffersnow.com/spiritonline/internal/iwmaster/protocol"
)

func handleGetBots(i *protocol.IWContext) error {

	file, err := os.ReadFile("public/iwmaster/botnames.txt")
	if err != nil {
		return protocol.IWError_InvalidBotsFile
	}

	if len(file) == 0 {
		return protocol.IWError_InvalidBotsFile
	}

	if len(i.CommandInfo.Data) != 0 {
		return protocol.IWError_InvalidCommand
	}

	return i.Send(protocol.IWCommandInfo{
		Command: protocol.IWCommand_BotsResponse,
		Data: []string{
			string(file),
		},
	})
}
