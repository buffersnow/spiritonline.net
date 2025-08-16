package gp

import (
	"encoding/base64"
	"iter"
	"strconv"
)

type GamespyData struct {
	data string
}

func (conv GamespyData) DataBlock() []byte {
	return []byte(conv.data)
}

func (conv GamespyData) String() string {
	return conv.data
}

func (conv GamespyData) Integer() (int, error) {
	return strconv.Atoi(conv.data)
}

func (conv GamespyData) Dictionary() iter.Seq[GamespyKV] {
	return DeserializeKVDict(conv.data)
}

func (conv GamespyData) Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(conv.data))
}
