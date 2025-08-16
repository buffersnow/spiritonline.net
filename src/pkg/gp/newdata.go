package gp

import (
	"encoding/base64"
	"slices"
	"strconv"
)

type GameSpyNewData struct{}

var NewData GameSpyNewData

func (GameSpyNewData) String(key string, value string) GamespyKV {
	return GamespyKV{key: key, value: GamespyData{data: value}}
}

func (GameSpyNewData) DataBlock(key string, value []byte) GamespyKV {
	convData := string(value)
	return GamespyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyNewData) Integer(key string, value int) GamespyKV {
	convData := strconv.Itoa(value)
	return GamespyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyNewData) Integer64(key string, value int64) GamespyKV {
	convData := strconv.FormatInt(value, 10)
	return GamespyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyNewData) Dictionary(key string, value []GamespyKV) GamespyKV {
	convData := SerializeKVDict(slices.Values(value))
	return GamespyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyNewData) Base64(key string, value []byte) GamespyKV {
	convData := base64.StdEncoding.EncodeToString(value)
	return GamespyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyNewData) Boolean(key string, value bool) GamespyKV {
	convData := ""
	if value {
		convData = "1"
	}

	return GamespyKV{key: key, value: GamespyData{data: convData}}
}
