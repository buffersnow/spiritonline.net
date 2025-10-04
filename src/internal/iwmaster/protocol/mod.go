package protocol

import (
	"strings"

	"buffersnow.com/spiritonline/pkg/log"
	"buffersnow.com/spiritonline/pkg/net"
)

type IWContext struct {
	CommandInfo IWCommandInfo
	Connection  *net.UdpPacket
	Log         log.LoggingFactory
}

type IWCommandInfo struct {
	Command string
	Data    []string
}

func (i IWContext) Send(wci IWCommandInfo) error {

	msgBytes := []byte{0xff, 0xff, 0xff, 0xff}
	msgBytes = append(msgBytes, []byte(wci.Command)...)

	if len(wci.Data) != 0 {
		msgBytes = append(msgBytes, ' ')
		msgBytes = append(msgBytes, []byte(strings.Join(wci.Data, " "))...)
	}

	msgBytes = append(msgBytes, []byte{0xff, 0xff, 0xff}...)

	return i.Connection.WriteBytes(msgBytes)
}

func (i IWContext) Error(err *IWError) error {

	return i.Send(IWCommandInfo{
		Command: IWCommand_Error,
		Data: []string{
			err.ErrorCode,
			err.Message,
		},
	})
}
