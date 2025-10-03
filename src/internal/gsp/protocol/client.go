package protocol

import (
	"fmt"

	"buffersnow.com/spiritonline/pkg/gp"
)

func (g GamespyContext) Send(gci gp.GameSpyCommandInfo) error {

	msg := []gp.GameSpyKV{}
	if gci.SubCommand != 0 {
		msg = append(msg, gp.Message.Integer(gci.Command, gci.SubCommand))
	} else {
		msg = append(msg, gp.Message.Empty(gci.Command))
	}

	msg = append(msg, gci.Data...)
	return g.SendRaw(msg)
}

func (g GamespyContext) SendRaw(kvs []gp.GameSpyKV) error {

	data := ""
	for _, datapair := range kvs {
		data += datapair.Depickle()
	}
	data += "\\final\\"

	return g.Connection.WriteText(data)
}

func (g GamespyContext) Error(err gp.GameSpyError) error {

	msg := []gp.GameSpyKV{
		gp.Message.Integer("err", err.ErrorCode),
		gp.Message.String("errmsg", fmt.Sprintf("%s.", err.Message)), //~ yes, really i'm too lazy to add dots to the error
	}

	if err.IsFatal {
		msg = append(msg, gp.Message.Boolean("fatal", true))
		defer g.Connection.Close()
	}

	gci := gp.GameSpyCommandInfo{
		Command: GPCommand_Error,
		Data:    msg,
	}

	g.Log.Error("Error", "%s! (Error Code: 0x%04x, Fatal: %t)", err.Message, err.ErrorCode, err.IsFatal)

	return g.Send(gci)
}
