package gp

import (
	"encoding/base64"

	"github.com/spf13/cast"
)

type GameSpyData struct {
	data string
}

func (g GameSpyData) DataBlock() []byte {
	return []byte(g.data)
}

func (g GameSpyData) String() string {
	return g.data
}

func (g GameSpyData) Integer() (int, error) {
	return cast.ToIntE(g.data)
}

func (g GameSpyData) Dictionary() []GameSpyKV {
	return PickleDictionary(g.data)
}

func (g GameSpyData) Base64() string {
	return base64.StdEncoding.EncodeToString([]byte(g.data))
}
