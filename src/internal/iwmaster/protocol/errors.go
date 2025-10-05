package protocol

import (
	"fmt"
	"strings"
)

type IWError struct {
	ErrorCode string
	Message   string
}

// Error implements the error interface.
func (g *IWError) Error() string {
	return fmt.Sprintf("IW Error: %s - %s", g.ErrorCode, strings.ReplaceAll(g.Message, "\n", " "))
}

var (
	IWError_InvalidLocation   = &IWError{ErrorCode: "MEMPHIS", Message: "Failed to locate required service."}
	IWError_HeartbeatBlocked  = &IWError{ErrorCode: "ECHELON", Message: "Server has been blacklist from iwmaster.\nPlease contact an admin to resolve this."}
	IWError_InvalidGameType   = &IWError{ErrorCode: "NATGRID", Message: "Unsupported game type sent to iwmaster backend."}
	IWError_InvalidProtocol   = &IWError{ErrorCode: "BULLRUN", Message: "Invalid protocol version sent to iwmaster backend."}
	IWError_InvalidCommand    = &IWError{ErrorCode: "PINWALE", Message: "Invalid command or invalid command format was sent to iwmaster backend."}
	IWError_ChallengeMismatch = &IWError{ErrorCode: "MAINWAY", Message: "Challenge mismatched between iwmaster backend and server."}
	IWError_ValidationError   = &IWError{ErrorCode: "SPECTRA", Message: "Data sent to iwmaster backend failed validation."}
	IWError_InvalidBotsFile   = &IWError{ErrorCode: "CIVINTL", Message: "Botnames could not be loaded by the server."}
	IWError_Reserved          = &IWError{ErrorCode: "SPYCORE", Message: "Reserved"}
)
