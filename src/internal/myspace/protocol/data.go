package protocol

import (
	"encoding/base64"
	"iter"
	"strconv"
)

type MySpaceData struct {
	data string
}

func (conv MySpaceData) DataBlock() []byte {
	return []byte(conv.data)
}

func (conv MySpaceData) String() string {
	return conv.data
}

func (conv MySpaceData) Integer() (int, error) {
	return strconv.Atoi(conv.data)
}

func (conv MySpaceData) Dictionary() iter.Seq[MySpaceKV] {
	return DeserializeKVDict(conv.data)
}

func (conv MySpaceData) Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(conv.data))
}
