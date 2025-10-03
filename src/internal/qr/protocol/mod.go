package protocol

import (
	"buffersnow.com/spiritonline/pkg/gp"
	"buffersnow.com/spiritonline/pkg/log"
)

type QR2CommandInfo struct {
	Command byte
	Log     log.LoggingFactory
}

type QR1CommandInfo struct {
	Command gp.GameSpyCommandInfo
	Log     log.LoggingFactory
}
