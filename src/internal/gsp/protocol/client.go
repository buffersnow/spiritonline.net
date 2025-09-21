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
