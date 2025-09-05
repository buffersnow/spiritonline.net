package gp

import (
	"encoding/base64"
	"slices"

	"github.com/spf13/cast"
)

type GameSpyMessageBuilder struct{}

var Message GameSpyMessageBuilder

func (GameSpyMessageBuilder) String(key string, value string) GameSpyKV {
	return GameSpyKV{key: key, value: GamespyData{data: value}}
}

func (GameSpyMessageBuilder) DataBlock(key string, value []byte) GameSpyKV {
	convData := string(value)
	return GameSpyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyMessageBuilder) Integer(key string, value int) GameSpyKV {
	convData := cast.ToString(value)
	return GameSpyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyMessageBuilder) Integer64(key string, value int64) GameSpyKV {
	convData := cast.ToString(value)
	return GameSpyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyMessageBuilder) Dictionary(key string, value []GameSpyKV) GameSpyKV {
	convData := DepickleDictionary(slices.Values(value))
	return GameSpyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyMessageBuilder) Base64(key string, value []byte) GameSpyKV {
	convData := base64.StdEncoding.EncodeToString(value)
	return GameSpyKV{key: key, value: GamespyData{data: convData}}
}

func (GameSpyMessageBuilder) Boolean(key string, value bool) GameSpyKV {
	convData := ""
	if value {
		convData = "1"
	}

	return GameSpyKV{key: key, value: GamespyData{data: convData}}
}
