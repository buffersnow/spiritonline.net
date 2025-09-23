package protocol

import "buffersnow.com/spiritonline/pkg/gp"

func (g GamespyContext) Send(kvs []gp.GameSpyKV) error {

	data := ""
	for _, datapair := range kvs {
		data += datapair.Serialize()
	}
	data += "\\final\\"

	return g.Connection.WriteText(data)
}

func (g GamespyContext) Error(err gp.GameSpyError) error {

	msg := []gp.GameSpyKV{
		gp.Message.Empty("error"),
		gp.Message.Integer("err", err.ErrorCode),
		gp.Message.String("errmsg", err.Message),
	}

	if err.IsFatal {
		msg = append(msg, gp.Message.Boolean("fatal", true))
	}

	return g.Send(msg)
}
