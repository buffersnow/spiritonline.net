package protocol

import "buffersnow.com/spiritonline/pkg/log"

type QR2Context struct {
	Command byte
	Log     log.LoggingFactory
}
