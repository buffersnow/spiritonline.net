package protocol

import (
	"encoding/base64"
	"slices"
	"strconv"
)

type MySpaceNewData struct{}

var NewData MySpaceNewData

func (MySpaceNewData) String(key string, value string) MySpaceKV {
	return MySpaceKV{key: key, value: MySpaceData{data: value}}
}

func (MySpaceNewData) DataBlock(key string, value []byte) MySpaceKV {
	convData := string(value)
	return MySpaceKV{key: key, value: MySpaceData{data: convData}}
}

func (MySpaceNewData) Integer(key string, value int) MySpaceKV {
	convData := strconv.Itoa(value)
	return MySpaceKV{key: key, value: MySpaceData{data: convData}}
}

func (MySpaceNewData) Integer64(key string, value int64) MySpaceKV {
	convData := strconv.FormatInt(value, 10)
	return MySpaceKV{key: key, value: MySpaceData{data: convData}}
}

func (MySpaceNewData) Dictionary(key string, value []MySpaceKV) MySpaceKV {
	convData := SerializeKVDict(slices.Values(value))
	return MySpaceKV{key: key, value: MySpaceData{data: convData}}
}

func (MySpaceNewData) Base64(key string, value []byte) MySpaceKV {
	convData := base64.StdEncoding.EncodeToString(value)
	return MySpaceKV{key: key, value: MySpaceData{data: convData}}
}

func (MySpaceNewData) Boolean(key string, value bool) MySpaceKV {
	convData := ""
	if value {
		convData = "1"
	}

	return MySpaceKV{key: key, value: MySpaceData{data: convData}}
}
